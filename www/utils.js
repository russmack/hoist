

function convertUnixTime(unix_timestamp) {
    var date = new Date(unix_timestamp*1000);
    // hours part from the timestamp
    var hours = date.getHours();
    // // minutes part from the timestamp
    var minutes = "0" + date.getMinutes();
    // // seconds part from the timestamp
    var seconds = "0" + date.getSeconds();

    // // will display time in 10:30:23 format
    var formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
    //return formattedTime;
    return date;
}

function Table() {
    this.table;
}
Table.prototype.create = function(tableElem) {
    this.table = tableElem;
};
Table.prototype.setHeader = function(v) {
    var th = this.table.tHead;
    if (th !== null) {
        this.table.deleteTHead();
    }
    th = this.table.createTHead();
    var row = th.insertRow(0);
    for (var i=0; i<v.length; i++) {
        var cell = row.insertCell(i);
        var txt = document.createTextNode(v[i]);
        cell.appendChild(txt);
    }
};
Table.prototype.addBody = function(v) {
    var newtbody = document.createElement('tbody');
    var tbody = this.table.tBodies;
    if (tbody !== null) {
        tbody[0].parentNode.replaceChild(newtbody, tbody[0]);
    }
    var b = this.table.tBodies[0];
    for (var r=0; r<v.length; r++) {
        var row = b.insertRow();
        for (prop in v[r]) {
            if( v[r].hasOwnProperty(prop) ) {
                var cell = row.insertCell();
                var txt = document.createTextNode(v[r][prop]);
                cell.appendChild(txt);
            }
        }
    }
};
Table.prototype.clearHeader = function() {
    if (this.table.tHead !== null) {
        this.table.deleteTHead();
    }
};
Table.prototype.clearBody = function() {
    //for (var i=0; i<this.table.rows
    var newtbody = document.createElement('tbody');
    var tbody = this.table.tBodies;
    if (tbody !== null) {
        tbody[0].parentNode.replaceChild(newtbody, tbody[0]);
    }
};
Table.prototype.clear = function() {
    this.clearHeader();
    this.clearBody();
};
