var xmlhttp = CreateXMLHttpRequest();
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

    console.log("e.shiftKey =" + e.shiftKey + " e.Key=" + e.key);
    if (e.ctrlKey === true && e.key != 'Control') {

        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown", true);
        xmlCtrlKeyhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xmlCtrlKeyhttp.send("key_code=" + 17);
    }
    else if (e.shiftKey === true && e.key != 'Shift') {
        console.log("shift is enter");
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown", true);
        xmlCtrlKeyhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xmlCtrlKeyhttp.send("key_code=" + 16);
    }
    else if (e.altKey === true && e.key != 'Alt') {
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown", true);
        xmlCtrlKeyhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xmlCtrlKeyhttp.send("key_code=" + 18);
    }
    else if (e.metaKey === true && e.key != 'OS' && e.key != 'Meta') {
        xmlCtrlKeyhttp.onreadystatechange = callhandle;
        xmlCtrlKeyhttp.open("POST", "/keydown", true);
        xmlCtrlKeyhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xmlCtrlKeyhttp.send("key_code=" + 91);
    }


    keynum = convertKeyCodeToFirefoxStandard(keynum)
    xmlhttp.onreadystatechange = callhandle;
    xmlhttp.open("POST", "/keydown", true);
    xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xmlhttp.send("key_code=" + keynum);
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

    keynum = convertKeyCodeToFirefoxStandard(keynum)
    xmlhttp.onreadystatechange = callhandle;
    xmlhttp.open("POST", "/keyup", true);
    xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xmlhttp.send("key_code=" + keynum);
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