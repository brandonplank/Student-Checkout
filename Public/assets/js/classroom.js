const successNotification = window.createNotification({});

// CSV Table fetcher
function getTable() {
    $.ajax({
        type: "GET",
        url: "/GetCSV/",
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    });
}

// Run this automatically on page load
getTable()

function sendTestContent(content) {
    var request = new XMLHttpRequest();
    console.log("Going to send", content)
    request.open("POST", "/id/" + btoa(content), true);
    request.send();
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

// QR Code scanner

function sendStatusToWebPage() {
    let parsedJson = JSON.parse(this.responseText)
    if(parsedJson.isOut) {
        successNotification({
            title: 'Signed back in',
            message: parsedJson.name + ' has signed back in'
        });
    } else {
        successNotification({
            title: 'Signed out',
            message: parsedJson.name + ' has signed out'
        });
    }

    let request = new XMLHttpRequest();
    request.open("POST", "/id/" + btoa(parsedJson.name));
    request.send();
}

let scanner = new Instascan.Scanner({ video: document.getElementById('preview') });
scanner.addListener('scan', function (content) {
    console.log(content);

    var request = new XMLHttpRequest()
    request.timeout = 5000;
    request.addEventListener("load", sendStatusToWebPage);
    request.open("POST", "/isOut/" + btoa(content));
    request.send();

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
0
console.log("Hi reader :) This is Brandon here(Class of 2022) congrats on clicking F12 or view page src :P\n\nThis project was made using a multitude of languages, here is the list\n\nHTML(not really a programming language)\nJavaScript\nGoLang\n\nPlease always be nice to Mrs. Hart, she is the best teacher ever to exist.\nTalk to you on the flip side.")