package controller

import "stncCms/app/utils/mock"

var (
	userApp    mock.UserAppInterface
	foodApp    mock.PostAppInterface
	fakeUpload mock.UploadFileInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	s  = InitUsers(&userApp, &fakeAuth, &fakeToken)                       //We use all mocked data here
	f  = InitPost(&foodApp, &userApp, &fakeUpload, &fakeAuth, &fakeToken) //We use all mocked data here
	au = NewAuthenticate(&userApp, &fakeAuth, &fakeToken)                 //We use all mocked data here

)
