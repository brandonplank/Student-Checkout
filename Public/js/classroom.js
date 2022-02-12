var xhr = new XMLHttpRequest();
const myNotification = window.createNotification({
});
let scanner = new Instascan.Scanner({ video: document.getElementById('preview') });
function getTable() {
    $.ajax({
        type: "GET",
        url: "/GetCSV/",
        success: function (data) {
            //$('body').append(arrayToTable(Papa.parse(data).data));
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    });
}
getTable()
scanner.addListener('scan', function (content) {
    console.log(content);
    xhr.open("POST", "/id/" + btoa(content), true);
    xhr.send();
    myNotification({
        title: 'Success',
        message: 'Read the QR code! You are free to go.'
    });
    getTable()
});
Instascan.Camera.getCameras().then(function (cameras) {
    if (cameras.length > 0) {
        scanner.start(cameras[0]);
    } else {
        console.error('No cameras found.');
    }
}).catch(function (e) {
    console.error(e);
});

function sendTestContent(content) {
    console.log("Going to send", content)
    xhr.open("POST", "/id/" + btoa(content), true);
    xhr.send();
}

function arrayToTable(tableData) {
    var table = $('<table></table>');
    $(tableData).each(function (i, rowData) {
        var row = $('<tr></tr>');
        $(rowData).each(function (j, cellData) {
            row.append($('<td>'+cellData+'</td>'));
        });
        table.append(row);
    });
    return table;
}