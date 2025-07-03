package forms

type BaseForm struct {
	FieldErrors   map[string]string `form:"-"`
	GenericErrors []string          `form:"-"`
}

func (f *BaseForm) AddFieldError(field, message string) {
	if f.FieldErrors == nil {
		f.FieldErrors = make(map[string]string)
	}
	f.FieldErrors[field] = message
}

func (f *BaseForm) AddGenericError(message string) {
	if f.GenericErrors == nil {
		f.GenericErrors = make([]string, 1)
	}
	f.GenericErrors = append(f.GenericErrors, message)
}

type SignupForm struct {
	BaseForm
	Name     string `form:"name" validate:"required,notblank,min=1,max=64"`
	Email    string `form:"email" validate:"required,notblank,email,min=1,max=64"`
	Password string `form:"password" validate:"required,notblank,min=8,max=64"`
}

type LoginForm struct {
	BaseForm
	Email    string `form:"email" validate:"required,notblank,email,min=1,max=64"`
	Password string `form:"password" validate:"required,notblank,min=8,max=64"`
}

type CreateHabitForm struct {
	BaseForm
	Name string `form:"name" validate:"required,notblank,min=1,max=32"`
}

type CreateEntryForm struct {
	BaseForm
	HabitID   int64  `form:"habitId" validate:"required,notblank"`
	EntryDate string `form:"entryDate" validate:"required,notblank"`
}
