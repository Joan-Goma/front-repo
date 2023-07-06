function getCookie() {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; neftAuth=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}
  // On load Toast
  setTimeout(function () {
    $.ajax({
          type: "GET",
          url: "https://test.joan-goma.repl.co/v1/secured/user",
          headers: {'neftAuth':getCookie()},
          success: function (data) {
            if (data.username !== "") {
              document.getElementById('userNameAPI').innerHTML = data.username;
              toastr['success'](
              'You have successfully logged in to Vuexy. Now you can start to explore!',
              '👋 Welcome ' +data.username +'!',
                {
                  closeButton: true,
                  tapToDismiss: false,
                }
              );
            } else {
              document.getElementById('userNameAPI').innerHTML = "Sr/a. Fantasma";
            }
            
          }
      });
    
  }, 2000);