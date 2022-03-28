// Scheduler
function runAtTime(method, hour, minute, second) {
    (function loop() {
        var now = new Date();
        if (now.getHours() === hour && now.getMinutes() === minute && now.getSeconds() === second) {
            method();
        }
        now = new Date();
        var delay = 60000 - (now % 60000);
        setTimeout(loop, delay);
    })();
}

// CSV Table fetcher
function getTable() {
    $.ajax({
        type: "GET",
        url: "/GetAdminCSV/",
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

function SearchStudent() {
    var name = document.getElementById('name').value
    console.log("Searching for", name)
    document.getElementById('name').value = "";
    $.ajax({
        type: "POST",
        url: "/search/" + btoa(name),
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

function addTeacher(name, email) {
    console.log("Adding", name)
    $.ajax({
        type: "POST",
        url: "/addTeacher",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "name": name,
            "email": email
        }),
        success: function (data) {
        }
    })
}

function removeTeacher(email) {
    console.log("Removing", email)
    $.ajax({
        type: "POST",
        url: "/removeTeacher",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "email": email,
        }),
        success: function (data) {
        }
    })
}

$('#name').keypress(function(e){
    if (e.which === 13) {
        SearchStudent()
    }
});

// Run this automatically on page load
getTable()

// Run every 10 seconds
setInterval(getTable, 1000 * 30);

function arrayToTable(tableData) {
    var table = $('<table class="table"></table>');
    $(tableData).each(function (i, rowData) {
        var row = $('<tr></tr>');
        $(rowData).each(function (j, cellData) {
            if(cellData.length >= 1) {
                row.append($('<td>'+cellData+'</td>'))
            }
        })
        table.append(row)
    })
    return table
}