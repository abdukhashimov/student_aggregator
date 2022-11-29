package handlers

import (
	"net/http"

	en_locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	en := en_locale.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, translator)
}

func validatorWrapper[T any](handler func(input T, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
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

		handler(*input, w, r)
	}
}
