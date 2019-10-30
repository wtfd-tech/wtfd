import dialogPolyfill from 'dialog-polyfill';
export default class LoginRegister {
showDialog(dlg: HTMLDialogElement) {
    dlg.showModal();
    registerBackdropClickHandler(dlg);
}

constructor(){
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
        
        const userDialogButton = document.getElementById("userdialogbutton");
        const userDialog = document.getElementById("userdialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(userDialog);

        userDialogButton.addEventListener("click", function () {
            showDialog(userDialog);
        });
    }


  const loginForm = document.getElementById("loginform");
  const loginError = document.getElementById("loginerror");
  const loginSubmit = document.getElementById("loginsubmit");
  const loginLoading = document.getElementById("loginloading");

  loginForm.addEventListener("submit", function(e){
    loginLoading.style.display = "block";
    e.preventDefault();
    fetch("/login", {body: new URLSearchParams(new FormData(loginForm)), method: 'post'})
    .then((resp)=> resp.text())
    .then((resp) => {
      if (resp.includes("Server Error")) {
        loginLoading.style.display = "none";
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
  const registerSubmit = document.getElementById("registersubmit");
  const registerLoading = document.getElementById("registerloading");

  registerForm.addEventListener("submit", function(e){
    registerLoading.style.display = "block";
    e.preventDefault();
    fetch("/register", {body: new URLSearchParams(new FormData(registerForm)), method: 'post'})
    .then((resp)=> resp.text())
    .then((resp) => {
      registerLoading.style.display = "none";
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

  const edituserForm = document.getElementById("edituserform");
  const edituserError = document.getElementById("editusererror");
  const edituserSubmit = document.getElementById("editusersubmit");
  const edituserLoading = document.getElementById("edituserloading");

  edituserForm.addEventListener("submit", function(e){
    edituserLoading.style.display = "block";
    e.preventDefault();
    fetch("/changepassword", {body: new URLSearchParams(new FormData(edituserForm)), method: 'post'})
    .then((resp)=> resp.text())
    .then((resp) => {
      edituserLoading.style.display = "none";
      if (resp.includes("Server Error")) {
        edituserError.innerHTML = resp;
        edituserSubmit.setAttribute("class", "button fail");
        setTimeout(()=>{
            edituserError.innerHTML = "";
            edituserSubmit.setAttribute("class", "button");
        },"2000")
      } else {
        location.href = "/";
      }

    });
    return false;

  });
}
}
