package apierrors

type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}

type ForbiddenError string

func (e ForbiddenError) Error() string {
	return string(e)
}
