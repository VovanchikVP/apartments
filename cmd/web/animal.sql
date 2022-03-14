CREATE TABLE IF NOT EXISTS animals(
id int AUTO_INCREMENT,
name varchar(50),
age int,
PRIMARY KEY(id));

INSERT INTO animals VALUES(1,'Hippo',10);
INSERT INTO animals VALUES(2,'Ele',20);

CREATE TABLE IF NOT EXISTS counters(
    type varchar(100),
    number varchar(100),
    verification_date date,
    apartment_id int,
    FOREIGN KEY (apartment_id) REFERENCES apartments(ROWID)
);

CREATE TABLE IF NOT EXISTS address(
    post_index int,
    city varchar(100),
    street varchar(100),
    house varchar(20),
    apartment varchar(20)
);

CREATE TABLE IF NOT EXISTS property_documents(
    type varchar(100),
    number varchar(100),
    date varchar(100)
);

CREATE TABLE IF NOT EXISTS apartments(
    address_id int,
    count_rooms int,
    property_document_id int,
    rent int2,
    FOREIGN KEY (address_id) REFERENCES address(ROWID),
    FOREIGN KEY (property_document_id) REFERENCES property_documents(ROWID)
);

CREATE TABLE IF NOT EXISTS indications(
	counter_id int,
	date date,
	data float,
	FOREIGN KEY (counter_id) REFERENCES counters(ROWID)
);

CREATE TABLE IF NOT EXISTS tariffs(
    counter_id int,
	set_date date,
	cost float,
	FOREIGN KEY (counter_id) REFERENCES counters(ROWID)
);

CREATE TABLE IF NOT EXISTS type_pyments(
    name varchar(500)
);

CREATE TABLE IF NOT EXISTS id_cards(
    type varchar(100),
	number varchar(100),
	issued date
);

CREATE TABLE IF NOT EXISTS persons(
    last_name varchar(100),
	first_name varchar(100),
	patronymic varchar(100),
	id_card_id int,
	phone varchar(20),
	address_id int,
	FOREIGN KEY (id_card_id) REFERENCES id_cards(ROWID),
	FOREIGN KEY (address_id) REFERENCES address(ROWID)
);

CREATE TABLE IF NOT EXISTS contracts_rent(
    number varchar(50),
	date date,
	employer_id int,
	landlord_id int,
	apartment_id int,
	date_start_rent date,
	date_end_rent date,
	date_apartment_transfer date,
	rental float,
	date_rental date,
	deposit float,
	transferred_amount float,
	payments_communal int2,
	payments_network int2,
	payments_electric int2,
	payments_heating int2,
	payments_cold_water int2,
	payments_hot_water int2,
	additional_terms varchar,
	file_contract varchar,
	FOREIGN KEY (employer_id) REFERENCES persons(ROWID),
	FOREIGN KEY (landlord_id) REFERENCES persons(ROWID),
	FOREIGN KEY (apartment_id) REFERENCES apartments(ROWID)
);

CREATE TABLE IF NOT EXISTS tenant(
    contract_rent_id int,
	person_id int,
	FOREIGN KEY (contract_rent_id) REFERENCES contracts_rent(ROWID),
	FOREIGN KEY (person_id) REFERENCES persons(ROWID)
);

CREATE TABLE IF NOT EXISTS payments(
    apartment_id int,
	cost float,
	admission int2,
	type_payment_id int,
	date date,
	FOREIGN KEY (apartment_id) REFERENCES apartments(ROWID),
	FOREIGN KEY (type_payment_id) REFERENCES type_pyments(ROWID)
);


-- INSERT INTO counters VALUES ('холодная вода', '937332812', '2023-11-08'), ('холодная вода', '937332597', '2023-11-08');
