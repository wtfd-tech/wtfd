const stateObj = {foo: "baar"};

let flagsubmiteventlistenerfunc = function () {
};
let solutioneventlistenerfunc = function () {
};
let flaginputeventlistenerfunc = function () {
};

function addChallEventListener(title, points) {
    let elem = document.getElementById(title);
    elem.addEventListener("click", function () {
        let detView = document.getElementById("detailview");
        let detDescription = document.getElementById("detaildescription");
        let detTitle = document.getElementById("detailtitle");
        let detPoints = document.getElementById("detailpoints");
        let flagsubmitbutton = document.getElementById("flagsubmitbutton");
        let solutionbutton = document.getElementById("solutionbutton");
        let solutiondiv = document.getElementById("solutiondiv");
        let solutioninnerdiv = document.getElementById("solutioninnerdiv");
        let flagInput = document.getElementById("flaginput");
        let msgBox = document.getElementById("flagsubmitmsg");
        let checkLoading = document.getElementById("checkloading");
        let challUri = document.getElementById("challuri");
        history.pushState(stateObj, "challenge " + title, title);
        flagsubmitbutton.removeEventListener("click", flagsubmiteventlistenerfunc);
        flagInput.removeEventListener("keypress", flaginputeventlistenerfunc);
        solutionbutton.removeEventListener("click", solutioneventlistenerfunc);

        flagsubmiteventlistenerfunc = function () {
            const data = new URLSearchParams();
            checkLoading.style.display = "block";
            data.append("flag", flagInput.value);
            data.append("challenge", title);
            console.log("hey, the flag is " + flagInput.value);
            fetch("/submitflag", {method: 'post', body: data})
                .then(resp => resp.text())
                .then((resp) => {
                    checkLoading.style.display = "none";
                    if (resp === "correct") {
                        location.href = "/";
                    } else {
                        flagsubmitbutton.setAttribute("class", "button fail");
                        setTimeout(() => {
                            flagsubmitbutton.setAttribute("class", "button");
                        }, 1000);
                        msgBox.innerHTML = resp;
                    }
                });
            flagInput.value = "";
        };
        flaginputeventlistenerfunc = function (e) {
            if (e.key === 'Enter') {
                flagsubmiteventlistenerfunc();
            }
        };
        solutioneventlistenerfunc = function () {
            solutioninnerdiv.innerHTML = "<i>Loading, please wait...</i>";
            fetch("/solutionview/" + title)
                .then(response => response.text())
                .then((response) => {
                    solutioninnerdiv.innerHTML = response;
                });


        };
        flagInput.addEventListener("keypress", flaginputeventlistenerfunc);
        flagsubmitbutton.addEventListener("click", flagsubmiteventlistenerfunc);
        solutionbutton.addEventListener("click", solutioneventlistenerfunc);


        detView.addEventListener("close", function () {
            flagInput.value = "";
            msgBox.innerHTML = "";
            solutioninnerdiv.innerHTML = "";
        });

        detDescription.innerHTML = "<i>Loading, please wait...</i>";
        detTitle.innerHTML = "LOADING";
        detPoints.innerHTML = "-";
        fetch("/detailview/" + title).then(resp => resp.text()).then(function (response) {
            detDescription.innerHTML = response;
            detTitle.innerHTML = title;
            detPoints.innerHTML = points;
        });
        if (elem.getAttribute("class").includes("completed")) {
            flagsubmitbutton.style.display = "none";
            flagInput.style.display = "none";
            solutionbutton.style.display = "";
            solutiondiv.style.display = "";
        } else {
            flagsubmitbutton.style.display = "";
            flagInput.style.display = "";
            solutionbutton.style.display = "none";
            solutiondiv.style.display = "none";
        }

        challUri.style.display = "none";
        fetch("/uriview/" + title).then(resp => resp.text()).then(function (response) {
            if(response != "") {
                challUri.style.display = "";
                challUri.href = response;
            }
        });

        showDialog(detView);
    });

}

function registerBackdropClickHandler(dlg) {
    Array.prototype.slice.call(document.getElementsByClassName("backdrop")).forEach(function (elem) {
        elem.addEventListener("click", function () {
            history.pushState(stateObj, "index", "/");
            dlg.close();
        });
    });
}

function absolute(a) {
    return Math.abs(a);
}

function signum(a) {
    if (a < 0) return -1.0;
    if (a > 0) return 1.0;
    return a;
}

const svgNS = "http://www.w3.org/2000/svg";

function drawPath(svg, path, startX, startY, endX, endY, drawFunction, nothinginbetween) {
    // get the path's stroke width (if one wanted to be  really precize, one could use half the stroke size)
    const stroke = parseFloat(path.getAttribute("stroke-width"));
    // check if the svg is big enough to draw the path, if not, set heigh/width
    if (svg.getAttribute("height") < (endY + stroke)) svg.setAttributeNS(null, "height", (endY + stroke));
    if (svg.getAttribute("width") < (startX + stroke)) svg.setAttributeNS(null, "width", (startX + stroke));
    if (svg.getAttribute("width") < (endX + stroke)) svg.setAttributeNS(null, "width", (endX + stroke));

    //var deltaX = (endX - startX) * 0.15;
    //var deltaY = (endY - startY) * 0.15;

    const deltaNum = 25;
    const deltaX = (endX === startX) ? 0 : deltaNum;
    const deltaY = (endY === startY) ? 0 : deltaNum;

    // for further calculations which ever is the shortest distance
    const delta = deltaY < absolute(deltaX) ? deltaY : absolute(deltaX);
    console.log("deltax: " + deltaX + ", deltay:" + deltaY + ", delta: " + delta);
    // set sweep-flag (counter/clock-wise)
    // if start element is closer to the left edge,
    // draw the first arc counter-clockwise, and the second one clock-wise
    let arc1 = 0;
    let arc2 = 1;
    if (startX > endX) {
        arc1 = 1;
        arc2 = 0;
    }
    // draw tha pipe-like path
    // 1. move a bit down, 2. arch,  3. move a bit to the right, 4.arch, 5. move down to the end
    if (!drawFunction) {
        path.setAttributeNS(null, "d", "M" + startX + " " + startY +
            " H" + (startX + delta) +
            " A" + delta + " " + delta + " 0 0 " + arc2 + " " + (startX + 2 * delta) + " " + (startY + delta) +
            " V" + (endY - delta) +
            // " A" + delta + " " +  delta + " 0 0 " + arc1 + " " + (startX + 3*delta*signum(deltaX)) + " " + (startY + 3*delta) +
            // " H" + (endX - 3*delta*signum(deltaX)) +
            // " A" + delta + " " +  delta + " 0 0 " + arc2 + " " + (endX-2*delta) + " " + (startY + 4*delta) +
            // " V" + (endY-1*delta) +
            " A" + delta + " " + delta + " 0 0 " + arc1 + " " + (startX + 3 * delta) + " " + endY +
            " H" + (endX)
        );
    } else {
        if (startY === endY) {
            //75 is half of grid-column-gap
            if (nothinginbetween) {
            path.setAttributeNS(null, "d", "M" +startX + " " + startY + " H" + endX);

            } else {
            const mid = 75;
            path.setAttributeNS(null, "d", "M" + startX + " " + startY +
                " H" + (startX + mid - deltaNum) +
                " A" + deltaNum + " " + deltaNum + " 0 0 " + arc1 + " " + (startX + mid) + " " + (startY - deltaNum) +
                " V" + (startY - 2 * deltaNum) +
                " A" + deltaNum + " " + deltaNum + " 0 0 " + arc2 + " " + (startX + mid + deltaNum) + " " + (startY - 3 * deltaNum) +
                " H" + (endX - mid - deltaNum) +
                " A" + deltaNum + " " + deltaNum + " 0 0 " + arc2 + " " + (endX - mid) + " " + (endY - 2 * deltaNum) +
                " V" + (endY - deltaNum) +
                " A" + deltaNum + " " + deltaNum + " 0 0 " + arc1 + " " + (endX - mid + deltaNum) + " " + (endY) +
                " H" + endX
            );
            }
        } else {
            path.setAttributeNS(null, "d", "M" + startX + " " + startY +
                " H" + (startX + delta) +
                " A" + delta + " " + delta + " 0 0 " + arc2 + " " + (startX + 2 * delta) + " " + (startY + delta) +
                " V" + (startY + 2 * delta) +
                " A" + delta + " " + delta + " 0 0 " + arc1 + " " + (startX + 3 * delta * signum(deltaX)) + " " + (startY + 3 * delta) +
                " H" + (endX - 3 * delta * signum(deltaX)) +
                " A" + delta + " " + delta + " 0 0 " + arc2 + " " + (endX - 2 * delta) + " " + (startY + 4 * delta) +
                " V" + (endY - delta) +
                " A" + delta + " " + delta + " 0 0 " + arc1 + " " + (endX - delta) + " " + endY +
                " H" + (endX)
            );
        }
    }
}

function connectElementss(svg, startElem, endElems, color) {
    let elem = document.getElementById(startElem);
    endElems.forEach(function (item) {
        console.log("start: " + startElem + " end: " + item);
        let selem = document.getElementById(item);
        connectElements(svg, elem, selem, color)

    });

}

function isInbetween(startElem){
vals = colnum.values();
    a = vals.next();
  let done = false
  while(!done){
    console.log(a.value.row, startElem.row, a.value.col, parseInt(startElem.col)+1);
    if(a.value.row === startElem.row && a.value.col === startElem.col+1) return true;
    a = vals.next();
    done = a.done;
  }
  return false;

}

function connectElements(svg, startElem, endElem, color) {
    const drawFunction = colnum.get(endElem.id).col - colnum.get(startElem.id).col > 1;
    
    const path = document.createElementNS(svgNS, "path");
    path.setAttributeNS(null, "d", "M0 0");
    path.setAttributeNS(null, "stroke", color);
    path.setAttributeNS(null, "fill", "none");
    path.setAttributeNS(null, "stroke-width", "12px");
    svg.insertBefore(path, svg.lastChild);
    const svgContainer = document.getElementById("svgContainer");

    // if first element is lower than the second, swap!
    if (startElem.offsetLeft > endElem.offsetLeft) {
        const temp = startElem;
        startElem = endElem;
        endElem = temp;
    }

    // get (top, left) corner coordinates of the svg container   
    //const svgTop = svgContainer.offsetTop; //Unused
    const svgLeft = svgContainer.offsetLeft;

    // get (top, left) coordinates for the two elements
    const startCoord = {left: startElem.offsetLeft, top: startElem.offsetTop};
    const endCoord = {left: endElem.offsetLeft, top: endElem.offsetTop};

    // calculate path's start (x,y)  coords
    // we want the x coordinate to visually result in the element's mid point
    const startX = startCoord.left + startElem.offsetWidth - svgLeft;    // x = left offset + 0.5*width - svg's left offset
    const startY = startCoord.top - 0.5 * (startElem.offsetHeight);        // y = top offset + height - svg's top offset

    // calculate path's end (x,y) coords
    const endX = endCoord.left;// + 0.5*endElem.offsetWidth - svgLeft;
    const endY = endCoord.top - 0.5 * (endElem.offsetHeight);

    // call function for drawing the path
    drawPath(svg, path, startX, startY, endX, endY, drawFunction, !isInbetween(colnum.get(startElem.id)));

}

(function () {
    let svg1 = document.getElementById("svg1");
    window.addEventListener("resize", function () {
        svg1.setAttribute("width", "0");
        svg1.setAttribute("height", "0");
        svg1.innerHTML = "";
        connectAll();
    });

    const detailview = document.getElementById("detailview");
    // noinspection JSUnresolvedVariable
    dialogPolyfill.registerDialog(detailview);
    const detailclosebutton = document.getElementById("detailclosebutton");
    detailclosebutton.addEventListener("click", function () {
        detailview.close();

        history.pushState(stateObj, "index", "/");

    });

    svg1.setAttribute("width", "0");
    svg1.setAttribute("height", "0");
    svg1.innerHTML = "";
    // noinspection JSUnresolvedVariable
    flagsubmitbutton.addEventListener("click", flagsubmiteventlistenerfunc);
    start();
})();
