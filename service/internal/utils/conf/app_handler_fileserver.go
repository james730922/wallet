package conf

type FileServerHandler struct{}

func (FileServerHandler) GetPath() string {
	return appConf.v.GetString("fileserver.path")
}

func (FileServerHandler) GetInternalPath() string {
	return appConf.v.GetString("fileserver.config.path")
}

func (FileServerHandler) GetInternalFolder() string {
	return appConf.v.GetString("fileserver.config.folder")
}

func (FileServerHandler) GetHost() string {
	return appConf.v.GetString("fileserver.config.host")
}

func (FileServerHandler) GetUser() string {
	return appConf.v.GetString("fileserver.config.user_name")
}

func (FileServerHandler) GetPassword() string {
	return appConf.v.GetString("fileserver.config.password")
}
