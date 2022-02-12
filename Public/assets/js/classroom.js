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
    setTimeout(function () {
        getTable()
    }, 500);
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
    var table = $('<table class="table"></table>');
    $(tableData).each(function (i, rowData) {
        var row = $('<tr></tr>');
        $(rowData).each(function (j, cellData) {
            if(cellData.length >= 1) {
                row.append($('<td>'+cellData+'</td>'));
            }
        });
        table.append(row);
    });
    return table;
}

console.log("Hi reader :) This is Brandon here(Class of 2022) congrats on clicking F12 or view page src :P\n\nThis project was made using a multitude of languages, here is the list\n\nHTML(not really a programming language)\nJavaScript\nGoLang\n\nPlease always be nice to Mrs. Hart, she is the best teacher ever to exist.\nTalk to you on the flip side.")