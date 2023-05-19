function addContractRent(data){
    let blank_row = $(".contract_rent-row");
    let table_contract = blank_row.parent();
    for (let i=0; i<data.length; i++){
        let new_row = blank_row.clone()
        let contract = data[i]
        let list_td = new_row.find('td')
        $(list_td[0]).html(contract['ID'])
        $(list_td[1]).html(contract['Number'])
        $(list_td[2]).html(getAddress(contract))
        $(list_td[3]).html(contract['DateRental'])
        $(list_td[4]).html(contract['Rental'])
        $(list_td[5]).html(contract['DateEndRent'])
        new_row.removeClass("hide")
        table_contract.append(new_row)
    }

}

function getAddress(data){
    let address = data['Apartment']['Address']
    return address['City'] + ' ' + address['Street'] + ' ' + address['House'] + '-' + ['Apartment']
}

function createSelect(select_id, url, func){
    $.ajax({
        type: 'get',
        url: url,
        success: function(data) {
            data = JSON.parse(data);
            for(let i=0; i<select_id.length; i++){
                let obj = $("#" + select_id[i])
                for(let j=0; j<data.length; j++){
                    func(obj, data[j])
                }
            }
        },
        error: function (jqXHR, exception) {
            console.log(jqXHR)
            console.log(exception)
        }
    });
}

function selectPersons(obj, data){
    obj.append($('<option>', {value: data['ID'],
        html: data['FirstName'] + ' ' + data['LastName'] + ' ' + data['Patronymic']}))
}

function selectApartment(obj, data){
    let address = data['Address']
    obj.append($('<option>', {value: data['ID'],
        html: address['City'] + ' ' + address['Street'] + ' ' + address['House'] + '-' + address['Apartment']}))
}

function selectIDCard(obj, data){
    obj.append($('<option>', {value: data['ID'],
        html: data['Type'] + ' ' + data['Number']}))
}


function selectAddress(obj, data){
    obj.append($('<option>', {value: data['ID'],
        html: data['City'] + ' ' + data['Street'] + ' ' + data['House'] + '-' + data['Apartment']}))
}

createSelect(['id_employer', 'id_landlord'], '/person?id=0&json=1', selectPersons)
createSelect(['id_apartment'], '/apartment?id=0&json=1', selectApartment)
createSelect(['id_id_card'], '/id_cards?id=0&json=1', selectIDCard)
createSelect(['id_address'], '/address?id=0&json=1', selectAddress)