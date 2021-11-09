$(document).ready(function () {
    but = [
        {
            extend: 'colvisGroup',
            text: 'DaiLy',
            show: [1, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21],
            hide: [2, 3, 4, 5, 6, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33]
        },
        {
            extend: 'colvisGroup',
            text: 'Other',
            show: [1, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33],
            hide: [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]
        },
        {
            extend: 'colvisGroup',
            text: 'Show all',
            show: ':hidden'
        }
    ]

    $("#datatableDay").DataTable({
        // bPaginate: true,
        // sPaginationType: "full_numbers",
        // bLengthChange: true,
        // bInfo: true,
        dom: 'Bfrtip',
        buttons: but,
        aLengthMenu: [10, 20, 50, 100]
    })
    $("#datatableTen").DataTable({
        dom: 'Bfrtip',
        buttons: but,
        aLengthMenu: [10,40]
    })
    $("#datatableMonth").DataTable({
        dom: 'Bfrtip',
        buttons: but,
        aLengthMenu: [12, 10]
    })
    $("#datatableQuarter").DataTable({
        dom: 'Bfrtip',
        buttons: but,
        aLengthMenu: [5, 10]
    })

    $("#datatableYear").DataTable({
        dom: 'Bfrtip',
        buttons: but,
        aLengthMenu: [5, 10]
    })

    $("#datatableNow").DataTable({
        dom: 'Bfrtip',
        buttons: but,
        aaSorting: false,
        aLengthMenu: [15, 15]
    })

    $("#datatableSummaryWeek").DataTable({
        aLengthMenu: [6, 12, 18],
        // bPaginate: false,
        aaSorting: false,
    })

    $("#datatableSummaryMonth").DataTable({
        aLengthMenu: [6, 12, 18],
        aaSorting: false,
    })

    $("#datatableSummaryQuarter").DataTable({
        aLengthMenu: [6, 12, 18],
        aaSorting: false,
    })

    $("#datatableSummaryYear").DataTable({
        aLengthMenu: [6, 12, 18],
        aaSorting: false,
    })

    $("#datatableRuns").DataTable({
        aLengthMenu: [12, 20, 50, 100],
    })

    $("#datatableField").DataTable({
        aLengthMenu: [10, 20, 50, 100],
        aaSorting: [[1, "asc"]],
    })

    $("#essentialMb").DataTable({
        dom: 'Bfrtip',
        aLengthMenu: [10, 20, 50, 100],
        buttons: [
            {
                extend: 'columnsToggle',
                // columns: '.toggle'
            },
            {
                extend: 'colvisGroup',
                text: 'Hide',
                show: [1],
                hide: [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33]
            }
        ]
    })

    // $("#datatable-buttons").DataTable({
    //     lengthChange: !1,
    //     buttons: ["copy", "excel", "pdf", "colvis"]
    // }).buttons().container().appendTo("#datatable-buttons_wrapper .col-md-6:eq(0)"), $(".dataTables_length select").addClass("form-select form-select-sm")
});
