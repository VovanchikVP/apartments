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