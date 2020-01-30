import './dragscroll';
import LoginRegister from './login-register';

(function(){

console.log("test");

new LoginRegister();

switch(location.pathname){
  case '/admin': {
    import( './leaderboard').then(_ => {
    new _.default();
    });
    import('./admin').then(_ => {
    new _.default();
    });
    break;
  }
  case '/leaderboard': {
    import('./leaderboard').then(_ => {
    new _.default();
    });
    break;
  }
  default: {
    import('./main').then(_ => {
    new _.default();
    });
    break;
  }
}

})();
