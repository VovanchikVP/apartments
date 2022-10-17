package entities

type Animal struct {
	ID   int
	Name string
	Age  int
}

type Counter struct {
	// Счетчик
	ID               int
	Type             string
	Number           string
	VerificationDate string
	Apartment        Apartment
}

type Indication struct {
	// Показание
	ID      int
	Counter Counter
	Date    string
	Data    float32
}

type Tariff struct {
	// Тариф
	ID      int
	Counter Counter
	SetDate string
	Cost    float32
}

type TypePayment struct {
	// Тип платежа ++
	ID   int
	Name string
}

type IDCard struct {
	// Документ удостоверяющий личность ++
	ID     int
	Type   string
	Number string
	Issued string
}

type Person struct {
	// Человек ++
	ID         int
	LastName   string
	FirstName  string
	Patronymic string
	IDCard     IDCard
	Phone      string
	Address    Address
}

type Address struct {
	// Адрес регистрации ++
	ID        int
	Index     int
	City      string
	Street    string
	House     string
	Apartment string
}

type PropertyDocuments struct {
	// Документ о собственности ++
	ID     int
	Type   string
	Number string
	Date   string
}

type Apartment struct {
	// Квартира
	ID                int
	Address           Address
	CountRooms        int
	PropertyDocuments PropertyDocuments
	Rent              bool
}

type ContractRent struct {
	// Договор аренды
	ID                    int
	Number                string
	Date                  string
	Employer              Person
	Landlord              Person
	Apartment             Apartment
	DateStartRent         string
	DateEndRent           string
	DateApartmentTransfer string
	Rental                float32
	DateRental            string
	Deposit               float32
	TransferredAmount     float32
	PaymentsCommunal      bool
	PaymentsNetwork       bool
	PaymentsElectric      bool
	PaymentsHeating       bool
	PaymentsColdWater     bool
	PaymentsHotWater      bool
	AdditionalTerms       string
	FileContract          string
}

type Tenant struct {
	ID           int
	ContractRent ContractRent
	Person       Person
}

type Payment struct {
	ID        int
	Apartment Apartment
	Cost      float32
	Admission bool
	Type      TypePayment
	Date      string
}
