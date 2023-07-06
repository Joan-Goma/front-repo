 /*=========================================================================================
    File Name: app-user-list.js
    Description: User List page
    --------------------------------------------------------------------------------------
    Item Name: Vuexy  - Vuejs, HTML & Laravel Admin Dashboard Template
    Author: PIXINVENT
    Author URL: http://www.themeforest.net/user/pixinvent

==========================================================================================*/
function getCookie() {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; neftAuth=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}
$(function () {
  var roleName = {
              1: "Soldado",
              2: "Traidor",
              3: "Lider equipo",
              4: "Lider peloton",
              5: "Comandante",
              6: "Organizador campo",
              7: "Admin",
            };
  ('use strict');

  var dtUserTable = $('.user-list-table'),
    newUserSidebar = $('.new-user-modal'),
    newUserForm = $('.add-new-user'),
    select = $('.select2'),
    dtContact = $('.dt-contact'),
    statusObj = {
      0: { title: 'Deshabilitado', class: 'badge-light-danger' },
      1: { title: 'Activado', class: 'badge-light-success' }
    };

  var userView = 'app-user-view-account.html';

  select.each(function () {
    var $this = $(this);
    $this.wrap('<div class="position-relative"></div>');
    $this.select2({
      // the following code is used to disable x-scrollbar when click in select input and
      // take 100% width in responsive also
      dropdownAutoWidth: true,
      width: '100%',
      dropdownParent: $this.parent()
    });
  });
  var dataUserApi
  // Users List datatable
  if (dtUserTable.length) {
    dtUserTable.DataTable({
      ajax: {
        type: "GET",
    dataType: "json",
    url: "https://APINEFT.joangoma.repl.co/api/secured/users",
    headers: {'neftAuth':getCookie()},
    dataSrc: '',
        complete: function (data) {
          $('#totalUsers').html(data.responseJSON.length)
        }
      
      }, // JSON file to add data
      columns: [
        // columns according to JSON
        { data: 'full_name' },
        { data: 'username' },
        { data: 'role' },
        { data: 'email' },
        { data: 'activated' },
        { data: 'photo' },
        { data: '' }
      ],
      columnDefs: [
        {
          // For Responsive
          className: 'control',
          orderable: false,
          responsivePriority: 2,
          targets: 0,
          render: function (data, type, full, meta) {
            return '';
          }
        },
        {
          // User full name and username
          targets: 1,
          responsivePriority: 4,
          render: function (data, type, full, meta) {
            var $username = full['username'],
              $full_name = full['full_name'],
              $image = full['photo'];
            if ($image) {
              // For Avatar image
              var $output =
                '<img src="' + assetPath + 'images/avatars/' + $image + '" alt="Avatar" height="32" width="32">';
            } else {
              // For Avatar badge
              var stateNum = Math.floor(Math.random() * 6) + 1;
              var states = ['success', 'danger', 'warning', 'info', 'dark', 'primary', 'secondary'];
              var $state = states[stateNum],
                $username = full['username'],
                $initials = $username.match(/\b\w/g) || [];
              $initials = (($initials.shift() || '') + ($initials.pop() || '')).toUpperCase();
              $output = '<span class="avatar-content">' + $initials + '</span>';
            }
            var colorClass = $image === '' ? ' bg-light-' + $state + ' ' : '';
            // Creates full output for row
            var $row_output =
              '<div class="d-flex justify-content-left align-items-center">' +
              '<div class="avatar-wrapper">' +
              '<div class="avatar ' +
              colorClass +
              ' me-1">' +
              $output +
              '</div>' +
              '</div>' +
              '<div class="d-flex flex-column">' +
              '<a href="' +
              userView +
              '" class="user_name text-truncate text-body"><span class="fw-bolder">' +
              $username +
              '</span></a>' +
              '<small class="emp_post text-muted">' +
              $full_name +
              '</small>' +
              '</div>' +
              '</div>';
            return $row_output;
          }
        },
        {
          // User Role
          targets: 2,
          render: function (data, type, full, meta) {
            var $role = full['role'];
            var roleBadgeObj = {
              1: feather.icons['circle'].toSvg({ class: 'font-medium-3 text-primary me-50' }),
              2: feather.icons['alert-triangle'].toSvg({ class: 'font-medium-3 text-warning me-50' }),
              3: feather.icons['chevron-up'].toSvg({ class: 'font-medium-3 text-success me-50' }),
              4: feather.icons['chevrons-up'].toSvg({ class: 'font-medium-3 text-info me-50' }),
              5: feather.icons['crosshair'].toSvg({ class: 'font-medium-3 text-danger me-50' }),
              6: feather.icons['tool'].toSvg({ class: 'font-medium-3 text-danger me-50' }),
              7: feather.icons['key'].toSvg({ class: 'font-medium-3 text-danger me-50' })
            };
            return "<span class='text-truncate align-middle'>" + roleBadgeObj[$role] + roleName[$role] + '</span>';
          }
        },
        {
          targets: 4,
          render: function (data, type, full, meta) {
            var $billing = full['billing'];

            return '<span class="text-nowrap">' + $billing + '</span>';
          }
        },
        {
          // User Status
          targets: 5,
          render: function (data, type, full, meta) {
            var $status = full['activated'];
            
            return (
              '<span class="badge rounded-pill ' +
              statusObj[$status].class +
              '" text-capitalized>' +
              statusObj[$status].title +
              '</span>'
            );
          }
        },
        {
          // Actions
          targets: -1,
          title: 'Actions',
          orderable: false,
          render: function (data, type, full, meta) {
            return (
              '<div class="btn-group">' +
              '<a class="btn btn-sm dropdown-toggle hide-arrow" data-bs-toggle="dropdown">' +
              feather.icons['more-vertical'].toSvg({ class: 'font-small-4' }) +
              '</a>' +
              '<div class="dropdown-menu dropdown-menu-end">' +
              '<a href="' +
              userView +
              '" class="dropdown-item">' +
              feather.icons['file-text'].toSvg({ class: 'font-small-4 me-50' }) +
              'Details</a>' +
              '<a href="javascript:;" class="dropdown-item delete-record">' +
              feather.icons['trash-2'].toSvg({ class: 'font-small-4 me-50' }) +
              'Delete</a></div>' +
              '</div>' +
              '</div>'
            );
          }
        }
      ],
      order: [[1, 'desc']],
      dom:
        '<"d-flex justify-content-between align-items-center header-actions mx-2 row mt-75"' +
        '<"col-sm-12 col-lg-4 d-flex justify-content-center justify-content-lg-start" l>' +
        '<"col-sm-12 col-lg-8 ps-xl-75 ps-0"<"dt-action-buttons d-flex align-items-center justify-content-center justify-content-lg-end flex-lg-nowrap flex-wrap"<"me-1"f>B>>' +
        '>t' +
        '<"d-flex justify-content-between mx-2 row mb-1"' +
        '<"col-sm-12 col-md-6"i>' +
        '<"col-sm-12 col-md-6"p>' +
        '>',
      language: {
        sLengthMenu: 'Show _MENU_',
        search: 'Buscar',
        searchPlaceholder: 'Buscar..'
      },
      // Buttons with Dropdown
      buttons: [
        {
          extend: 'collection',
          className: 'btn btn-outline-secondary dropdown-toggle me-2',
          text: feather.icons['external-link'].toSvg({ class: 'font-small-4 me-50' }) + 'Export',
          buttons: [
            {
              extend: 'print',
              text: feather.icons['printer'].toSvg({ class: 'font-small-4 me-50' }) + 'Print',
              className: 'dropdown-item',
              exportOptions: { columns: [1, 2, 3, 4, 5] }
            },
            {
              extend: 'csv',
              text: feather.icons['file-text'].toSvg({ class: 'font-small-4 me-50' }) + 'Csv',
              className: 'dropdown-item',
              exportOptions: { columns: [1, 2, 3, 4, 5] }
            },
            {
              extend: 'excel',
              text: feather.icons['file'].toSvg({ class: 'font-small-4 me-50' }) + 'Excel',
              className: 'dropdown-item',
              exportOptions: { columns: [1, 2, 3, 4, 5] }
            },
            {
              extend: 'pdf',
              text: feather.icons['clipboard'].toSvg({ class: 'font-small-4 me-50' }) + 'Pdf',
              className: 'dropdown-item',
              exportOptions: { columns: [1, 2, 3, 4, 5] }
            },
            {
              extend: 'copy',
              text: feather.icons['copy'].toSvg({ class: 'font-small-4 me-50' }) + 'Copy',
              className: 'dropdown-item',
              exportOptions: { columns: [1, 2, 3, 4, 5] }
            }
          ],
          init: function (api, node, config) {
            $(node).removeClass('btn-secondary');
            $(node).parent().removeClass('btn-group');
            setTimeout(function () {
              $(node).closest('.dt-buttons').removeClass('btn-group').addClass('d-inline-flex mt-50');
            }, 50);
          }
        },
        {
          text: 'Add New User',
          className: 'add-new btn btn-primary',
          attr: {
            'data-bs-toggle': 'modal',
            'data-bs-target': '#modals-slide-in'
          },
          init: function (api, node, config) {
            $(node).removeClass('btn-secondary');
          }
        }
      ],
      // For responsive popup
      responsive: {
        details: {
          display: $.fn.dataTable.Responsive.display.modal({
            header: function (row) {
              var data = row.data();
              return 'Detalles de ' + data['name'];
            }
          }),
          type: 'column',
          renderer: function (api, rowIdx, columns) {
            var data = $.map(columns, function (col, i) {
              return col.columnIndex !== 6 // ? Do not show row in modal popup if title is blank (for check box)
                ? '<tr data-dt-row="' +
                    col.rowIdx +
                    '" data-dt-column="' +
                    col.columnIndex +
                    '">' +
                    '<td>' +
                    col.title +
                    ':' +
                    '</td> ' +
                    '<td>' +
                    col.data +
                    '</td>' +
                    '</tr>'
                : '';
            }).join('');
            return data ? $('<table class="table"/>').append('<tbody>' + data + '</tbody>') : false;
          }
        }
      },
      language: {
        paginate: {
          // remove previous & next text from pagination
          previous: '&nbsp;',
          next: '&nbsp;'
        }
      },
      initComplete: function () {
        // Adding role filter once table initialized
        this.api()
          .columns(2)
          .every(function () {
            var column = this;
            var label = $('<label class="form-label" for="UserRole">Rol</label>').appendTo('.user_role');
            var select = $(
              '<select id="UserRole" class="form-select text-capitalize mb-md-0 mb-2"><option value=""> Seleccionar rol </option></select>'
            )
              .appendTo('.user_role')
              .on('change', function () {
                var val = $.fn.dataTable.util.escapeRegex($(this).val());
                column.search(val ? '^' + val + '$' : '', true, false).draw();
              });

            column
              .data()
              .unique()
              .sort()
              .each(function (d, j) {
                select.append('<option value="' + roleName[d] + '" class="text-capitalize">' + roleName[d] + '</option>');
              });
          });
        // Adding status filter once table initialized
        this.api()
          .columns(5)
          .every(function () {
            var column = this;
            var label = $('<label class="form-label" for="FilterTransaction">Estado</label>').appendTo('.user_status');
            var select = $(
              '<select id="FilterTransaction" class="form-select text-capitalize mb-md-0 mb-2"><option value=""> Selecciona un estado </option></select>'
            )
              .appendTo('.user_status')
              .on('change', function () {
                var val = $.fn.dataTable.util.escapeRegex($(this).val());
                column.search(val ? '^' + val + '$' : '', true, false).draw();
              });

            column
              .data()
              .unique()
              .sort()
              .each(function (d, j) {
                select.append(
                  '<option value="' +
                    statusObj[1].title +
                    '" class="text-capitalize">' +
                    statusObj[1].title +
                    '</option>' +
                  '<option value="' +
                    statusObj[0].title +
                    '" class="text-capitalize">' +
                    statusObj[0].title +
                    '</option>'
                );
              });
          });
      }
    });
  }

  // Form Validation
  newUserForm.on('submit', function (e) {
    e.preventDefault();
    var data = JSON.stringify(detectNumeric(getFormData(newUserForm.serializeArray())))
      $.ajax({
          type: "POST",
          dataType: 'json',
          data: data,
          url: "https://APINEFT.joangoma.repl.co/api/secured/users",
          headers: {'neftAuth':getCookie()},
          statusCode: {
            201: handle201,
            500: handle500
          }
      });
      e.preventDefault();
    });
});

var dtUserTable = $('.user-list-table'),
    newUserSidebar = $('.new-user-modal')
var handle201 = function(data) {
            toastr['success']('Congrats new user created!', 'Success!', {
              closeButton: true,
              tapToDismiss: false,
            });
            newUserSidebar.modal('hide');
            var oTable = dtUserTable.DataTable();
            // to reload
            oTable.ajax.reload();
};

var handle500 = function(data) {
  newUserSidebar.modal('hide');
            var jsonData = JSON.stringify(data)
            Swal.fire({
              icon: 'error',
              title: 'Oops...',
              text: 'Something went wrong!',
              footer: data.responseJSON.error
            })
}
function detectNumeric(obj) {
  for (var index in obj) {
    // if object property value *is* a number, like 1 or "500"
    if (obj[index] === ''){
    }
    else if (!isNaN(obj[index])) {
      // convert it to 1 or 500
      obj[index] = Number(obj[index]);
    }
  }
  
  return obj
}
//utility function
function getFormData(data) {
   var unindexed_array = data;
   var indexed_array = {};

   $.map(unindexed_array, function(n, i) {
    indexed_array[n['name']] = n['value'];
   });

   return indexed_array;
}