package masterunitlist

type LabelValue struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Group string `json:"-"`
}

type quickListResponse struct {
	Units []Unit `json:"Units"`
}

type Unit struct {
	ID        int     `json:"Id"`
	Name      string  `json:"Name"`
	Class     string  `json:"Class"`
	Variant   string  `json:"Variant"`
	Tonnage   float64 `json:"Tonnage"`
	Rules     string  `json:"Rules"`
	TRO       string  `json:"TRO"`
	DateIntro string  `json:"DateIntroduced"`
	ImageURL  string  `json:"ImageUrl"`

	Technology struct {
		ID   int    `json:"Id"`
		Name string `json:"Name"`
	} `json:"Technology"`

	Type struct {
		ID   int    `json:"Id"`
		Name string `json:"Name"`
	} `json:"Type"`

	Role struct {
		ID   int    `json:"Id"`
		Name string `json:"Name"`
	} `json:"Role"`

	BFType             string `json:"BFType"`
	BFSize             int    `json:"BFSize"`
	BFMove             string `json:"BFMove"`
	BFTMM              int    `json:"BFTMM"`
	BFArmor            int    `json:"BFArmor"`
	BFStructure        int    `json:"BFStructure"`
	BFDamageShort      int    `json:"BFDamageShort"`
	BFDamageShortMin   bool   `json:"BFDamageShortMin"`
	BFDamageMedium     int    `json:"BFDamageMedium"`
	BFDamageMediumMin  bool   `json:"BFDamageMediumMin"`
	BFDamageLong       int    `json:"BFDamageLong"`
	BFDamageLongMin    bool   `json:"BFDamageLongMin"`
	BFDamageExtreme    int    `json:"BFDamageExtreme"`
	BFDamageExtremeMin bool   `json:"BFDamageExtemeMin"`
	BFOverheat         int    `json:"BFOverheat"`
	BFPointValue       int    `json:"BFPointValue"`
	BFAbilities        string `json:"BFAbilities"`
}

type Era struct {
	ID   int
	Name string
}
