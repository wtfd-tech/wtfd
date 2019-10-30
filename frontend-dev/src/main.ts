import dialogPolyfill from 'dialog-polyfill'

export default class MainPage {

stateObj = {foo: "baar"};

flagsubmiteventlistenerfunc = function () {
};
solutioneventlistenerfunc = function () {
};
flaginputeventlistenerfunc = function () {
};

// Used for selection dialog in bugreport
selCategory = document.getElementById("bugreportcategory");

addChallEventListener(title: string, points: number) {
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
        let challAuthor = document.getElementById("challauthor");
        history.pushState(this.stateObj, "challenge " + title, title);
        flagsubmitbutton.removeEventListener("click", this.flagsubmiteventlistenerfunc);
        flagInput.removeEventListener("keypress", this.flaginputeventlistenerfunc);
        solutionbutton.removeEventListener("click", this.solutioneventlistenerfunc);

        this.flagsubmiteventlistenerfunc = function () {
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
        challAuthor.innerHTML = "LOADING..."
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

        fetch("/authorview/" + title).then(resp => resp.text()).then(function (response) {
            if(response != "") {
                challAuthor.innerHTML = response;
            } else {
                challAuthor.innerHTML = "&lt;Unknown&gt;"
            }
        });

        selCategory.value = title;
        showDialog(detView);
    });

}

registerBackdropClickHandler(dlg: HTMLDialogElement) {
    Array.prototype.slice.call(document.getElementsByClassName("backdrop")).forEach(function (elem: Element) {
        elem.addEventListener("click", function () {
            history.pushState(stateObj, "index", "/");
            dlg.close();
        });
    });
}

const absolute(a: number): number {
    return Math.abs(a);
}

const signum(a: number): number {
    if (a < 0) return -1.0;
    if (a > 0) return 1.0;
    return a;
}

const svgNS = "http://www.w3.org/2000/svg";

drawPath(svg: Element, path: SVGElement, startX: number, startY: number, endX: number, endY: number, drawFunction: boolean, nothinginbetween: boolean) {
    // get the path's stroke width (if one wanted to be  really precize, one could use half the stroke size)
    const stroke = parseFloat(path.getAttribute("stroke-width"));
    // check if the svg is big enough to draw the path, if not, set heigh/width
    if (Number(svg.getAttribute("height")) < (endY + stroke)) svg.setAttributeNS(null, "height", String(endY + stroke));
    if (Number(svg.getAttribute("width")) < (startX + stroke)) svg.setAttributeNS(null, "width", String(startX + stroke));
    if (Number(svg.getAttribute("width")) < (endX + stroke)) svg.setAttributeNS(null, "width", String(endX + stroke));

    //var deltaX = (endX - startX) * 0.15;
    //var deltaY = (endY - startY) * 0.15;

    const deltaNum = 25;
    const deltaX = (endX === startX) ? 0 : deltaNum;
    const deltaY = (endY === startY) ? 0 : deltaNum;

    // for further calculations which ever is the shortest distance
    const delta = deltaY < this.absolute(deltaX) ? deltaY : this.absolute(deltaX);
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
                " A" + delta + " " + delta + " 0 0 " + arc1 + " " + (startX + 3 * delta * this.signum(deltaX)) + " " + (startY + 3 * delta) +
                " H" + (endX - 3 * delta * this.signum(deltaX)) +
                " A" + delta + " " + delta + " 0 0 " + arc2 + " " + (endX - 2 * delta) + " " + (startY + 4 * delta) +
                " V" + (endY - delta) +
                " A" + delta + " " + delta + " 0 0 " + arc1 + " " + (endX - delta) + " " + endY +
                " H" + (endX)
            );
        }
    }
}

connectElementss(svg: Element, startElem: string, endElems: string[], color: string) {
    let elem = document.getElementById(startElem);
    endElems.forEach(function (item) {
        console.log("start: " + startElem + " end: " + item);
        let selem = document.getElementById(item);
        connectElements(svg, elem, selem, color)

    });

}

isInbetween(startElem: any){
  return new Promise(function(resolve){
let vals: any  = colnum.values();
    let a: any = vals.next();
  while(!a.done){
    //console.log(a.value, startElem, a.value.col, parseInt(startElem.col)+1, a.value.row === startElem.row && parseInt(a.value.col)-1 === parseInt(startElem.col));
    if(a.value.row === startElem.row && parseInt(a.value.col)-1 === parseInt(startElem.col)) resolve(true);
    a = vals.next();
    if(a.done){
  resolve(false);
    }
  }
  });

}

connectElements(svg: HTMLElement, startElem: HTMLElement, endElem: HTMLElement, color: string) {
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

  console.log(startElem);
    isInbetween(colnum.get(startElem.id)).then((ibt) => {
    drawPath(svg, path, startX, startY, endX, endY, drawFunction, !ibt);

    });
    // call function for drawing the path

}

constructer() {
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

        history.pushState(this.stateObj, "index", "/");

    });

    svg1.setAttribute("width", "0");
    svg1.setAttribute("height", "0");
    svg1.innerHTML = "";
    // noinspection JSUnresolvedVariable

    // Submit bugreport form
    flagsubmiteventlistenerfunc = function () {
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
    }

    // BUGREPORT STUFF
    flagsubmitbutton.addEventListener("click", flagsubmiteventlistenerfunc);
    let btnBugreport = document.getElementById("bugreport");
    let btnBugreportMain = document.getElementById("mainbugreport");
    let dlgBugreport = document.getElementById("bugreportview");
    let btnBugreportClose = document.getElementById("bugreportclosebutton");
    let btnBugreportSubmit = document.getElementById("bugreportbutton");
    let txtBugreportSubject = document.getElementById("subjectinput");
    let txtBugreportContent = document.getElementById("contentinput");
    let bugreportCheckLoading = document.getElementById("bugloading");
    btnBugreportClose.addEventListener("click", function () {
        dlgBugreport.close();
    });
    dialogPolyfill.registerDialog(dlgBugreport);
    btnBugreport.addEventListener("click", function() {
        showDialog(dlgBugreport);
    });
    btnBugreportMain.addEventListener("click", function() {
        selCategory.value = "Main Page";
        showDialog(dlgBugreport);
    });
    btnBugreportSubmit.addEventListener("click", function () {
        const data = new URLSearchParams();
        bugreportCheckLoading.style.display = "block";
        data.append("subject", "[" + selCategory.value + "] "
                    + txtBugreportSubject.value);
        data.append("content", txtBugreportContent.value);

        fetch("/reportbug", {method: 'POST', body: data})
            .then(resp => resp.text())
            .then((resp) => {
                bugreportCheckLoading.style.display = "none";
                if (resp === "OK") {
                    txtBugreportContent.value = "";
                    txtBugreportSubject.value = "";
                } else {
                    btnBugreportSubmit.setAttribute("class", "button fail bugreportsubmit");
                    setTimeout(() => {
                        btnBugreportSubmit.setAttribute("class", "button bug bugreportsubmit");
                    }, 1000);
                }
                bugreportCheckLoading.style.display = "none";
            });
    });

    // Add categories to bugreport selection
    bugreportCategories.forEach(function(elem) {
        var opt = document.createElement("option");
        opt.value= elem;
        opt.innerHTML = elem;
        selCategory.appendChild(opt);
    });
    //////// END BUGREPORT STUFF

    start();
}
}
