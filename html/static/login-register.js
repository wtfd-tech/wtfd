function showDialog(dlg) {
    dlg.showModal();
    registerBackdropClickHandler(dlg);
}

(function(){
    const loginButton = document.getElementById("loginbutton");
    if (loginButton != null) {
        const loginDialogCancelButton = document.getElementById("logincancelbutton");
        const loginDialog = document.getElementById("logindialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(loginDialog);

        loginButton.addEventListener("click", function () {
            showDialog(loginDialog);
        });
        loginDialogCancelButton.addEventListener("click", function () {
            loginDialog.close();
        });

        const registerButton = document.getElementById("registerbutton");
        const registerDialogCancelButton = document.getElementById(
            "registercancelbutton"
        );
        const registerDialog = document.getElementById("registerdialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(registerDialog);

        const loginTabButton = document.getElementById("logintabbutton");

        const loginUserBox = document.getElementById("loginuserbox");
        const loginPassBox = document.getElementById("loginpassbox");
        const registerUserBox = document.getElementById("registeruserbox");
        const registerPassBox = document.getElementById("registerpassbox");

        registerButton.addEventListener("click", function () {
            registerUserBox.value = loginUserBox.value;
            registerPassBox.value = loginPassBox.value;
            loginDialog.close();
            showDialog(registerDialog);
        });
        registerDialogCancelButton.addEventListener("click", function () {
            registerDialog.close();
        });

        loginTabButton.addEventListener("click", function () {
            loginUserBox.value = registerUserBox.value;
            loginPassBox.value = registerPassBox.value;
            registerDialog.close();
            showDialog(loginDialog);
        });
    } else {
        const logoutButton = document.getElementById("logoutbutton");
        logoutButton.addEventListener("click", function () {
            location.href = "/logout";
        });
    }


  const loginForm = document.getElementById("loginform");
  const loginError = document.getElementById("loginerror");
  const loginSubmit = document.getElementById("loginsubmit");
  loginForm.addEventListener("submit", function(e){
    e.preventDefault();
    fetch("/login", {body: new URLSearchParams(new FormData(loginForm)), method: 'post'})
    .then((resp)=> resp.text())
    .then((resp) => {
      if (resp.includes("Server Error")) {
        loginError.innerHTML = resp;
        loginSubmit.setAttribute("class", "button fail");
        setTimeout(()=>{
            loginError.innerHTML = "";
            loginSubmit.setAttribute("class", "button");
        },2000);
      } else {
        location.href = "/";
      }

    });
    return false;

  });
  const registerForm = document.getElementById("registerform");
  const registerError = document.getElementById("registererror");
  const registerSubmit = document.getElementById("registersubmit")
  registerForm.addEventListener("submit", function(e){
    e.preventDefault();
    console.log(new URLSearchParams(new FormData(registerForm)))
    fetch("/register", {body: new URLSearchParams(new FormData(registerForm)), method: 'post'})
    .then((resp)=> resp.text())
    .then((resp) => {
      if (resp.includes("Server Error")) {
        registerError.innerHTML = resp;
        registerSubmit.setAttribute("class", "button fail");
        setTimeout(()=>{
            registerError.innerHTML = "";
            registerSubmit.setAttribute("class", "button");
        },"2000")
      } else {
        location.href = "/";
      }

    });
    return false;

  });
})();
