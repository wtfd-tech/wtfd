function absolute(a){
  return Math.abs(a);
}
function signum(a){
if (a<0) return -1.0;
if (a>0) return 1.0;
  return a;
}
const svgNS = "http://www.w3.org/2000/svg";  
function drawPath(svg, path, startX, startY, endX, endY) {
    // get the path's stroke width (if one wanted to be  really precize, one could use half the stroke size)
    var stroke =  parseFloat(path.getAttribute("stroke-width"));
    // check if the svg is big enough to draw the path, if not, set heigh/width
    if (svg.getAttribute("height") <  (endY+stroke))                 svg.setAttributeNS(null, "height", (endY+stroke));
    if (svg.getAttribute("width" ) < (startX + stroke) )    svg.setAttributeNS(null, "width", (startX + stroke));
    if (svg.getAttribute("width" ) < (endX   + stroke) )    svg.setAttributeNS(null, "width", (endX   + stroke));
    
    var deltaX = (endX - startX) * 0.15;
    var deltaY = (endY - startY) * 0.15;
    // for further calculations which ever is the shortest distance
    var delta  =  deltaY < absolute(deltaX) ? deltaY : absolute(deltaX);

    // set sweep-flag (counter/clock-wise)
    // if start element is closer to the left edge,
    // draw the first arc counter-clockwise, and the second one clock-wise
    var arc1 = 0; var arc2 = 1;
    if (startX > endX) {
        arc1 = 1;
        arc2 = 0;
    }
    // draw tha pipe-like path
    // 1. move a bit down, 2. arch,  3. move a bit to the right, 4.arch, 5. move down to the end 
    path.setAttributeNS(null, "d",  "M"  + startX + " " + startY +
                    " H" + (startX + delta) +
                    " A" + delta + " " +  delta + " 0 0 " + arc2 + " " + (startX+2*delta) + " " + (startY + 1*delta) +
                    " V" + (startY + 2*delta) +
                    " A" + delta + " " +  delta + " 0 0 " + arc1 + " " + (startX + 3*delta*signum(deltaX)) + " " + (startY + 3*delta) +
                    " H" + (endX - 3*delta*signum(deltaX)) + 
                    " A" + delta + " " +  delta + " 0 0 " + arc2 + " " + (endX-2*delta) + " " + (startY + 4*delta) +
                    " V" + (endY-1*delta) +
                    " A" + delta + " " +  delta + " 0 0 " + arc1 + " " + (endX-1*delta) + " " + (endY - 0*delta) +
                    " H" + (endX) );
}

function connectElementss(svg, startElems, endElem,color){
  elem = document.getElementById(endElem);
  startElems.forEach(function(item){
    selem = document.getElementById(item)
    connectElements(svg, selem, elem,color)

  });

}

function connectElements(svg, startElem, endElem, color) {
    var path = document.createElementNS(svgNS,"path");
    path.setAttributeNS(null,"d","M0 0");
    path.setAttributeNS(null,"stroke",color);
    path.setAttributeNS(null,"fill","none");
    path.setAttributeNS(null,"stroke-width","12px");
    svg.insertBefore(path,svg.lastChild);
    var svgContainer= document.getElementById("svgContainer");

    // if first element is lower than the second, swap!
    if(startElem.offsetLeft > endElem.offsetLeft){
        var temp = startElem;
        startElem = endElem;
        endElem = temp;
    }

    // get (top, left) corner coordinates of the svg container   
    var svgTop  = svgContainer.offsetTop;
    var svgLeft = svgContainer.offsetLeft;

    // get (top, left) coordinates for the two elements
    var startCoord = {left: startElem.offsetLeft, top: startElem.offsetTop};
    var endCoord   = {left: endElem.offsetLeft, top: endElem.offsetTop};

    // calculate path's start (x,y)  coords
    // we want the x coordinate to visually result in the element's mid point
    var startX = startCoord.left + startElem.offsetWidth - svgLeft;    // x = left offset + 0.5*width - svg's left offset
    var startY = startCoord.top- 0.5*(startElem.offsetHeight);        // y = top offset + height - svg's top offset

        // calculate path's end (x,y) coords
    var endX = endCoord.left;// + 0.5*endElem.offsetWidth - svgLeft;
    var endY = endCoord.top  - 0.5*(endElem.offsetHeight);
    // call function for drawing the path
    drawPath(svg, path, startX, startY, endX, endY);

}

(function() {
  var loginButton = document.getElementById("loginbutton");
  if (loginButton != null) {
    var loginDialogCancelButton = document.getElementById("logincancelbutton");
    var loginDialog = document.getElementById("logindialog");
    dialogPolyfill.registerDialog(loginDialog);

    loginButton.addEventListener("click", function() {
      loginDialog.showModal();
    });
    loginDialogCancelButton.addEventListener("click", function() {
      loginDialog.close();
    });

    var registerButton = document.getElementById("registerbutton");
    var registerDialogCancelButton = document.getElementById(
      "registercancelbutton"
    );
    var registerDialog = document.getElementById("registerdialog");
    dialogPolyfill.registerDialog(registerDialog);

    registerButton.addEventListener("click", function() {
      loginDialog.close();
      registerDialog.showModal();
    });
    registerDialogCancelButton.addEventListener("click", function() {
      registerDialog.close();
    });
  } else {
    var logoutButton = document.getElementById("logoutbutton");
    logoutButton.addEventListener("click", function() {
      location.href = "/logout";
    });
  }
  
  svg1 = document.getElementById("svg1")
  window.addEventListener("resize", function(){
    svg1.setAttribute("width","0");
    svg1.setAttribute("height","0");
    svg1.innerHTML= "";
    connectAll();
  });
  svg1.setAttribute("width","0");
  svg1.setAttribute("height","0");
    svg1.innerHTML= "";
  connectAll();
})();
