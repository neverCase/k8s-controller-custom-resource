package proto

// The shortage ot the harbor projects
type HarborHub struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

type HarborHubList struct {
	Items []HarborHub `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborProject struct {
	ProjectID int32  `json:"projectId" protobuf:"varint,1,opt,name=projectId"`
	Name      string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type HarborProjectList struct {
	Items []HarborProject `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborRepository struct {
	RepositoryID int32  `json:"repositoryId" protobuf:"bytes,1,opt,name=repositoryId"`
	Name         string `json:"name" protobuf:"bytes,2,opt,name=name"`
	ProjectID    int32  `json:"projectId" protobuf:"varint,3,opt,name=projectId"`
}

type HarborRepositoryList struct {
	Items []HarborRepository `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborTag struct {
	Digest string `json:"digest" protobuf:"bytes,1,rep,name=digest"`
	Name   string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type HarborTagList struct {
	Items []HarborTag `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborRequest struct {
	HarborUrl string        `json:"harborUrl" protobuf:"bytes,1,opt,name=harborUrl"`
	Command   HarborCommand `json:"command" protobuf:"bytes,2,opt,name=command"`
	ProjectID int32         `json:"projectId" protobuf:"varint,3,opt,name=projectId"`
	ImageName string        `json:"imageName" protobuf:"bytes,4,opt,name=imageName"`
}
