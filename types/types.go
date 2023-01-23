package types

type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
	Owner   Owner
}

type Bucket struct {
	CreationDate string
	Name         string
}

type Owner struct {
	DisplayName string
	ID          string
}
