$(function () {
    const body = $('body');

    // Отправка формы добавления типа платежа.
    body.on('submit', '#id_form-type_payment', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        let url = '/type_payments';
        rest_ajax('post', data_form, url);
    });
    // Удаление типа платежа
    body.on('click', 'tr.type_payment-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/type_payments';
        let data = {'type_payment_id': id_row}
        rest_ajax('delete', data, url);
    });

    // Отправка формы добавления документа удостоверяющего личность.
    body.on('submit', '#id_form-id_card', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        console.log(data_form)
        let url = '/id_cards';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.id_card-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/id_cards';
        let data = {'id_card_id': id_row}
        rest_ajax('delete', data, url);
    });

    // Отправка формы добавления адреса.
    body.on('submit', '#id_form-address', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        console.log(data_form)
        let url = '/address';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.address-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/address';
        let data = {'address_id': id_row}
        rest_ajax('delete', data, url);
    });

    // Отправка формы добавления адреса.
    body.on('submit', '#id_form-property_document', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        console.log(data_form)
        let url = '/property_document';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.property_document-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/property_document';
        let data = {'property_document_id': id_row}
        rest_ajax('delete', data, url);
    });
    // Отправка формы добавления вида операции.
    body.on('submit', '#id_form-operation_groups', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        let url = '/operation_groups';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.operation_groups-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/operation_groups';
        let data = {'operation_groups_id': id_row}
        rest_ajax('delete', data, url);
    });
    // Отправка формы добавления операции.
    body.on('submit', '#id_form-operation', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        console.log(data_form)
        let url = '/operation';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.operation-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/operation';
        let data = {'operation_groups_id': id_row}
        rest_ajax('delete', data, url);
    });
    // Отправка формы добавления счетчика.
    body.on('submit', '#id_form-counter', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        console.log(data_form)
        let url = '/counter';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.counter-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/counter';
        let data = {'counter_id': id_row}
        rest_ajax('delete', data, url);
    });
    // Отправка формы добавления показаний.
    body.on('submit', '#id_form-indication', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        let url = '/indication';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.indication-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/indication';
        let data = {'indication_id': id_row}
        rest_ajax('delete', data, url);
    });
    // Отправка формы добавления тарифов.
    body.on('submit', '#id_form-tariff', function (event) {
        event.preventDefault();
        let data_form = $(this).serialize();
        let url = '/tariff';
        rest_ajax('post', data_form, url);
    });
    // Удаление
    body.on('click', 'tr.tariff-row td.tools.remove', function () {
        let id_row = $(this).closest('tr').attr('row_id');
        let url = '/tariff';
        let data = {'tariff_id': id_row}
        rest_ajax('delete', data, url);
    });
});

function rest_ajax(method, data, url) {
    $.ajax({
        type: method,
        url: url,
        data: data,
        success: function(data) {
            location.reload();
        },
        error: function (jqXHR, exception) {
            console.log(jqXHR)
            console.log(exception)
        }
    });
    return false;
}
