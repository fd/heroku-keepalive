function AppViewModel() {
  var self = this;

  self.apps = ko.observableArray([]);

  self.load_apps = function() {
    $.getJSON("/apps", function(data){
      for (var i in data) {
        data[i].all_domains = [data[i].domain_name].concat(data[i].domains || []);
        data[i].all_domains.sort(function(l, r){
          return (l.domain < r.domain ? -1 : (l.domain == r.domain ? 0 : 1));
        });
      }
      self.apps(data);
      self.apps.sort(function(l, r){
        return (l.name < r.name ? -1 : (l.name == r.name ? 0 : 1));
      });

    });
  }

  setInterval(function(){ self.load_apps() }, 2000);
}

$(function(){
  window.App = new AppViewModel();
  ko.applyBindings(window.App);
  window.App.load_apps();
});
