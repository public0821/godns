var dnserverControllers = angular.module('dnserverControllers', []);
var dnserverApp = angular.module('dnserverApp', ['ngRoute', 'ui.bootstrap', 'dnserverControllers']);
//var dnserverApp = angular.module('dnserverApp', ['ngRoute', 'ui.bootstrap']);

dnserverApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/record', {
        templateUrl: 'record/manager.html',
        controller: 'RecordCtrl'
      }).
      when('/configure/basic', {
        templateUrl: 'configure/basic.html',
        controller: 'ConfigureBasicCtrl'
      }).
      when('/user/changepwd', {
        templateUrl: 'user/changepwd.html',
        controller: 'UserChangePwdCtrl'
      }).
      otherwise({
        redirectTo: ''
      });
  }]);

var RecordEditCtrl = function ($scope, $http, $modalInstance, record, global_scope) {
    if(record==null){
        $scope.record_type=1;
        $scope.record_ttl=0;
        $scope.record_name="";
        $scope.record_value="";
        $scope.record_id=0;
    }else{
        $scope.record_type=record.Type;
        $scope.record_ttl=record.Ttl;
        $scope.record_name=record.Name;
        $scope.record_value=record.Value;
        $scope.record_id=record.Id;
    }
  $scope.edit_record = function () {
        var record = {
            Id:this.record_id,
            Name: this.record_name,
            Type: parseInt(this.record_type),
            Value: this.record_value,
            Ttl: parseInt(this.record_ttl)
        };
      console.log(record);
        if(record.Name.length == 0 || record.Value.length == 0){
            $scope.record_tip='name and value field should be specified';
            return
        }
        if(isNaN(record.Ttl) || isNaN(record.Type)){
            $scope.record_tip='type and ttl should be number';
            return
        }
        if(record.Ttl<0){
            $scope.record_tip="ttl can't be negative";
            return
        }
    $http.post('/record/',record).success(function(data) {
      console.log(data);
        if(record.Id!=0){
            for (var i=0; i < global_scope.records.length; i++) {
                if(global_scope.records[i].Id == record.Id){
                    global_scope.records[i]=record
                    break;
                }
            };
        }
        $scope.record_tip='success';
    }).error(function(data, status, headers, config) {
        $scope.record_tip='server error';
    });
        
  };

  $scope.cancel = function () {
    $modalInstance.close();
    //$modalInstance.dismiss('cancel');
  };
};
//var RecordAddCtrl = function ($scope, $http, $modalInstance) {
    //$scope.record_type=1
    //$scope.record_ttl=60
  //$scope.add_record = function () {
        //var record = {
            //name: this.add_record_name,
            //type: parseInt(this.add_record_type),
            //value: this.add_record_value,
            //ttl: parseInt(this.add_record_ttl)
        //};
      //console.log(record);
        //if(isNaN(record.ttl) || isNaN(record.type)){
        //$scope.add_record_tip='ttl should be number';
        //return
        //}
    //$http.post('/record/',record).success(function(data) {
      //console.log(data);
        //$scope.add_record_tip='add success';
    //});
        
  //};

  //$scope.cancel = function () {
    //$modalInstance.close();
    ////$modalInstance.dismiss('cancel');
  //};
//};

dnserverControllers.controller('ConfigureBasicCtrl', ['$scope', '$http',
  function ($scope, $http, $modal) {
    $http.get('/sysoption/mode').success(function(data) {
      $scope.mode = data.Value;
      if($scope.mode == "forward"){
        $http.get('/forwardserver/').success(function(data) {
            $scope.fservers=data
        });
      }
      console.log(data);
    }).error(function(data, status, headers, config) {
        $scope.tip="server error";
        $scope.mode_change=function(){};
    });

    $scope.mode_change= function () {
      if($scope.mode == "forward"){
        $http.post('/sysoption/',{"name":"mode", "value":$scope.mode}).success(function(data) {
            $scope.tip="";
            $http.get('/forwardserver/').success(function(data) {
                $scope.fservers=data
            });
        }).error(function(data, status, headers, config) {
            $scope.tip="server error";
        });
      }
      else if($scope.mode == "recursion"){
        $http.post('/sysoption/',{"name":"mode", "value":$scope.mode}).success(function(data) {
                $scope.tip="";
        }).error(function(data, status, headers, config) {
            $scope.tip="server error";
        });
      }
  };
    $scope.delete_fserver= function (fserver) {
        $http.delete('/forwardserver/'+fserver.Ip).success(function(data) {
            for (var i=0; i < $scope.fservers.length; i++) {
                if($scope.fservers[i].Ip == fserver.Ip){
                    $scope.fservers.splice(i, 1)
                    break;
                }
            };
            $scope.tip="";
        }).error(function(data, status, headers, config) {
            $scope.tip="server error";
        });
  };
    $scope.add_fserver= function () {
        $http.post('/forwardserver/', {"Ip":$scope.fserver_ip}).success(function(data) {
            $scope.tip="";
            $scope.fservers.push({"Ip":$scope.fserver_ip});
        }).error(function(data, status, headers, config) {
            $scope.tip="server error";
        });
  };
}]);

dnserverControllers.controller('UserChangePwdCtrl', ['$scope', '$http',
  function ($scope, $http, $modal) {
      $scope.changePassword= function () {
            var params={
                "old_user":$scope.old_user,
                "old_password":$scope.old_password,
                "new_user":$scope.new_user,
                "new_password":$scope.new_password,
            };
            $http.post('/user/chpassword/',params).success(function(data) {
                $scope.tip='success';
            }).error(function(data, status, headers, config) {
                $scope.tip=data;
            });
      };
}]);

dnserverControllers.controller('RecordCtrl', ['$scope', '$http', '$modal',
  function ($scope, $http, $modal) {
      $scope.records=[];
    $http.get('/record').success(function(data) {
      $scope.records = data;
      console.log(data);
    });
  $scope.open_add_record = function () {
    var modalInstance = $modal.open({
      templateUrl: 'record_edit.html',
      controller: RecordEditCtrl,
        resolve: {
        record: function () {
          return null;
        },
        global_scope: function () {
          return $scope;
        }
      }
    });

  };
  $scope.open_edit_record = function (record) {
    var modalInstance = $modal.open({
      templateUrl: 'record_edit.html',
      controller: RecordEditCtrl,
        resolve: {
        record: function () {
          return record;
        },
        global_scope: function () {
          return $scope;
        }
      }
    });

  };
  $scope.delete_record = function (record) {
      console.log(record)
    $http.delete('/record/'+record.Id).success(function(data) {
        for (var i=0; i < $scope.records.length; i++) {
            if($scope.records[i].Id == record.Id){
                $scope.records.splice(i, 1)
                break;
            }
        };
    }).error(function(data, status, headers, config) {
      //
    });
  };
 
  $scope.query_name=""
  $scope.query_type=1
  $scope.query_value=""
  $scope.query = function () {
    $http.get('/record?name='+$scope.query_name+"&type="+$scope.query_type+"&value="+$scope.query_value)
        .success(function(data) {
      $scope.records = data;
      console.log(data);
    });
  };
}]);

dnserverControllers.controller('RecordAddCtrl', ['$scope', '$http',
  function ($scope, $http) {
 
  }]);
//phonecatApp.controller('PhoneListCtrl', function ($scope) {
  //$scope.phones = [
    //{'name': 'Nexus S',
     //'snippet': 'Fast just got faster with Nexus S.'},
    //{'name': 'Motorola XOOM™ with Wi-Fi',
     //'snippet': 'The Next, Next Generation tablet.'},
    //{'name': 'MOTOROLA XOOM™',
     //'snippet': 'The Next, Next Generation tablet.'}
  //];
//});
