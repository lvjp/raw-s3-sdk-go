package types

type LocationConstraint struct {
	LocationConstraint string `xml:",chardata"`
}

type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
	Owner   *Owner
}

type Bucket struct {
	CreationDate *string
	Name         *string
}

type Owner struct {
	DisplayName *string
	ID          *string
}
