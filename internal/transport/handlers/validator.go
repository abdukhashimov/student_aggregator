package handlers

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	enlocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	en := enlocale.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")
	validate = validator.New()
	err := entranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		panic(err)
	}
}

func validatorWrapper[T any](handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := new(T)

		if err := readJSON(r.Body, input); err != nil {
			sendUnprocessableEntityError(w, err)
			return
		}

		err := validate.Struct(input)

		if err != nil {
			var errs []string
			for _, v := range err.(validator.ValidationErrors).Translate(translator) {
				errs = append(errs, v)
			}
			sendValidationError(w, errs)
			return
		}

		r = inputToContext(r, input)

		handler(w, r)
	}
}
