'use strict';
var app = angular.module('hotel_app', ['ngRoute']);
var BASE_URL = "http://127.0.0.1:8000/";

// Route providers
app.config(function ($routeProvider) {
  $routeProvider

    .when('/', {
      templateUrl: 'partials/home.html',
      controller: 'mainController',
      controllerAs: 'mainCtrl'
    })

    .when('/recent', {
      templateUrl: 'partials/recent.html',
      controller: 'recentController',
      controllerAs: 'recentCtrl'
    })

    .when('/search-result', {
      templateUrl: 'partials/result.html',
      controller: 'resultController',
      controllerAs: 'resultCtrl'
    })

    .otherwise({
      redirectTo: '/404'
    })

});


app.controller("mainController", function ($scope, $http, $location) {
  $scope.searchText = '';
  this.submitForm = function () {
    console.log($scope.searchText);
    console.log("Hello");
    $http({
      method: 'POST',
      data: {
        "query": $scope.searchText
      },
      url: BASE_URL + 'search',
      headers: {
        'Content-Type': 'application/json'
      },
    }).then(function successCallback(response) {
        console.log(response);
        $location.path('/recent');
      },
      function errorCallback(response) {
        console.log("Error");
        console.log(response);
      });
  }
});

app.controller("recentController", function ($scope, $http) {
  $scope.recentQueries = [];
  // this.submitForm = function () {
  // console.log($scope.searchText);
  console.log("Hello");
  $http({
    method: 'GET',
    data: {},
    url: BASE_URL + 'recent',
    headers: {
      'Content-Type': 'application/json'
    },
  }).then(function successCallback(response) {
      console.log(response.data);
      $scope.recentQueries = response.data;
    },
    function errorCallback(response) {
      console.log("Error");
      console.log(response);
    });
  // }

  this.getRowClass = function (status) {
    if (status === 'Completed') {
      return "table-success";
    } else if (status === 'Failed') {
      return "table-danger";
    } else {
      return "table-primary";
    }
  };
});

app.controller("resultController", function ($scope, $http, $location) {
  $scope.result = {};
  $http({
    method: 'GET',
    data: {},
    url: BASE_URL + 'search-result?queryId=' + $location.search().queryId,
    headers: {
      'Content-Type': 'application/json'
    },
  }).then(function successCallback(response) {
      console.log(response.data);
      $scope.result = response.data;
    },
    function errorCallback(response) {
      console.log("Error");
      console.log(response);
    });

});