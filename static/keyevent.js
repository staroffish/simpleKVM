var xmlCtrlKeyhttp = CreateXMLHttpRequest();
var isFirefox = isFirefoxBrowser();

function CreateXMLHttpRequest() {
    if (window.ActiveXObject) {
        xmlReq = new ActiveXObject("Microsoft.XMLHTTP");
    }
    else if (window.XMLHttpRequest) {
        xmlReq = new XMLHttpRequest();
    }
    return xmlReq;
}
function callhandle() { }
function keyDown(e) {
    if (window.event) // IE
    {
        keynum = e.keyCode
    }
    else if (e.which) // Netscape/Firefox/Opera
    {
        keynum = e.which
    }

    if (e.ctrlKey === true && e.key != 'Control') {

        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown?key_code=" + 17, true);
        xmlCtrlKeyhttp.send();
    }
    else if (e.shiftKey === true && e.key != 'Shift') {
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown?key_code=" + 16, true);
        xmlCtrlKeyhttp.send();
    }
    else if (e.altKey === true && e.key != 'Alt') {
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown?key_code=" + 18, true);
        xmlCtrlKeyhttp.send();
    }
    else if (e.metaKey === true && e.key != 'OS' && e.key != 'Meta') {
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown?key_code=" + 91, true);
        xmlCtrlKeyhttp.send();
    }

    var xmlDownHttp = CreateXMLHttpRequest();
    keynum = convertKeyCodeToFirefoxStandard(keynum)
    xmlDownHttp.onreadystatechange = callhandle;
    xmlDownHttp.open("POST", "/keydown?key_code=" + keynum, true);
    xmlDownHttp.send();
    return false;
}

function keyUp(e) {
    if (window.event) // IE
    {
        keynum = e.keyCode
    }
    else if (e.which) // Netscape/Firefox/Opera
    {
        keynum = e.which
    }
    var xmlUpHttp = CreateXMLHttpRequest();
    keynum = convertKeyCodeToFirefoxStandard(keynum)
    xmlUpHttp.onreadystatechange = callhandle;
    xmlUpHttp.open("POST", "/keyup?key_code=" + keynum, true);
    xmlUpHttp.send();
    return false;
}

function mouseMove(e) {
    var xmlmouseHttp = CreateXMLHttpRequest();
    xmlmouseHttp.onreadystatechange = callhandle;
    xmlmouseHttp.open("POST", "/mousemove?x=" + e.clientX + "&y=" + e.clientY, true);
    xmlmouseHttp.send();
    return false;
}

function mouseDown(e) {
    var xmlmouseHttp = CreateXMLHttpRequest();
    xmlmouseHttp.onreadystatechange = callhandle;
    xmlmouseHttp.open("POST", "/mousedown?button=" + e.button, true);
    xmlmouseHttp.send();
    return false;
}

function mouseUp(e) {
    var xmlmouseHttp = CreateXMLHttpRequest();
    xmlmouseHttp.onreadystatechange = callhandle;
    xmlmouseHttp.open("POST", "/mouseup?button=" + e.button, true);
    xmlmouseHttp.send();
    return false;
}

function callhandle() {
    // if (xmlhttp.status != 200) {
    //     console.log(xmlhttp.responseText);
    // }
}

function isFirefoxBrowser() {
    var ua = navigator.userAgent.toLowerCase();
    return ua.match(/firefox\/([\d.]+)/) != null;
}

function convertKeyCodeToFirefoxStandard(keynum) {
    if (!isFirefox) {
        if (keynum == 0xBD) {
            keynum = 0xAD
        }
        else if (keynum == 0xBB) {
            keynum = 0x3D
        }
        else if (keynum == 0xBA) {
            keynum = 0x3B
        }
    }
    return keynum
}