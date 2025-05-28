package forms

type BaseForm struct {
	Errors map[string]string `form:"-"`
}

func (f *BaseForm) AddError(field, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[field] = message
}

type SignupForm struct {
	BaseForm
	Name     string `form:"name" validate:"required,notblank,min=1,max=64"`
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
