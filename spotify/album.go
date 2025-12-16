package spotify

func getAlbum(href string, auth *auth) (*album, error) {
	return getModel[album](href, auth)
}
