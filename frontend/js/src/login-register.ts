import dialogPolyfill from 'dialog-polyfill';
import {showDialog} from './util';
export default class LoginRegister {

constructor(){
    const loginButton = document.getElementById("loginbutton");
    if (loginButton != null) {
        const loginDialogCancelButton = <HTMLButtonElement> document.getElementById("logincancelbutton");
        const loginDialog = <HTMLDialogElement> document.getElementById("logindialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(loginDialog);

        loginButton.addEventListener("click", () => {
            showDialog(loginDialog);
        });
        loginDialogCancelButton.addEventListener("click", function () {
            loginDialog.close();
        });

        const registerButton = <HTMLButtonElement> document.getElementById("registerbutton");
        const registerDialogCancelButton = <HTMLButtonElement> document.getElementById(
            "registercancelbutton"
        );
        const registerDialog = <HTMLDialogElement> document.getElementById("registerdialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(registerDialog);

        const loginTabButton = document.getElementById("logintabbutton");

        const loginUserBox = <HTMLInputElement> document.getElementById("loginuserbox");
        const loginPassBox = <HTMLInputElement> document.getElementById("loginpassbox");
        const registerUserBox = <HTMLInputElement> document.getElementById("registeruserbox");
        const registerPassBox = <HTMLInputElement> document.getElementById("registerpassbox");

        registerButton.addEventListener("click", () => {
            registerUserBox.value = loginUserBox.value;
            registerPassBox.value = loginPassBox.value;
            loginDialog.close();
            showDialog(registerDialog);
        });
        registerDialogCancelButton.addEventListener("click", () => {
            registerDialog.close();
        });

        loginTabButton.addEventListener("click", ()=> {
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
        
        const userDialogButton = <HTMLButtonElement> document.getElementById("userdialogbutton");
        const userDialog = <HTMLDialogElement> document.getElementById("userdialog");
        // noinspection JSUnresolvedVariable
        dialogPolyfill.registerDialog(userDialog);

        userDialogButton.addEventListener("click", ()=> {
            showDialog(userDialog);
        });
    }


  const loginForm = <HTMLFormElement> document.getElementById("loginform");
  const loginError = document.getElementById("loginerror");
  const loginSubmit = document.getElementById("loginsubmit");
  const loginLoading = document.getElementById("loginloading");

  loginForm.addEventListener("submit", function(e){
    loginLoading.style.display = "block";
    e.preventDefault();
    // @ts-ignore
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
  
  const registerForm = <HTMLFormElement> document.getElementById("registerform");
  const registerError = document.getElementById("registererror");
  const registerSubmit = document.getElementById("registersubmit");
  const registerLoading = document.getElementById("registerloading");

  registerForm.addEventListener("submit", (e) => {
    e.preventDefault();
    registerLoading.style.display = "block";
    // @ts-ignore
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
        },2000)
      } else {
        location.href = "/";
      }

    });
    e.preventDefault();
    e.stopPropagation();

  });

  const edituserForm = <HTMLFormElement> document.getElementById("edituserform");
  const edituserError = document.getElementById("editusererror");
  const edituserSubmit = document.getElementById("editusersubmit");
  const edituserLoading = document.getElementById("edituserloading");

  edituserForm.addEventListener("submit", function(e){
    edituserLoading.style.display = "block";
    e.preventDefault();
    // @ts-ignore
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
        },2000)
      } else {
        location.href = "/";
      }

    });
    return false;

  });
}

}
