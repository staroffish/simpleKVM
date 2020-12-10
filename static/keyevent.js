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

function shortcut(arr) {

    for (var i = 0; i < arr.length; i++) {
        keyDown(arr[i]);
    }
    // for (var i = 0; i < arr.length; i++) {
    //     keyUp(arr[i]);
    // }
}

function keyDown(key) {
    var xmlDownHttp = CreateXMLHttpRequest();
    var keynum = convertKeyCodeToFirefoxStandard(key)
    xmlDownHttp.onreadystatechange = callhandle;
    xmlDownHttp.open("POST", "/keydown?key_code=" + keynum, true);
    xmlDownHttp.send();
}

function keyUp(key) {
    var xmlDownHttp = CreateXMLHttpRequest();
    var keynum = convertKeyCodeToFirefoxStandard(key)
    xmlDownHttp.onreadystatechange = callhandle;
    xmlDownHttp.open("POST", "/keyup?key_code=" + keynum, true);
    xmlDownHttp.send();
}

function onKeyDown(e) {
    if (window.event) // IE
    {
        keynum = e.keyCode
    }
    else if (e.which) // Netscape/Firefox/Opera
    {
        keynum = e.which
    }

    if (e.ctrlKey === true && e.key != 'Control') {
        keyDown(17)
    }
    else if (e.shiftKey === true && e.key != 'Shift') {
        keyDown(16)
    }
    else if (e.altKey === true && e.key != 'Alt') {
        keyDown(18)
    }
    else if (e.metaKey === true && e.key != 'OS' && e.key != 'Meta') {
        keyDown(91)
    }

    keyDown(keynum)
    return false;
}

function onKeyUp(e) {
    if (window.event) // IE
    {
        keynum = e.keyCode
    }
    else if (e.which) // Netscape/Firefox/Opera
    {
        keynum = e.which
    }
    keyUp(keynum)
    return false;
}

function mouseMove(e) {
    var xmlmouseHttp = CreateXMLHttpRequest();
    xmlmouseHttp.onreadystatechange = callhandle;
    xmlmouseHttp.open("POST", "/mousemove?x=" + e.offsetX + "&y=" + e.offsetY, true);
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

function mouseScroll(e) {
    var xmlmouseHttp = CreateXMLHttpRequest();
    xmlmouseHttp.onreadystatechange = callhandle;
    xmlmouseHttp.open("POST", "/mousescroll?scroll=" + e.deltaY, true);
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