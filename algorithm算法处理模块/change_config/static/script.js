$(document).ready(function () {
  $("#add-container").click(function () {
    $("#containers").append(
      `<div class="form-inline mb-2">
              <input type="text" class="form-control mr-2" name="container_name" placeholder="Container Name">
              <input type="text" class="form-control" name="" placeholder="IP Address">
              <button type="button" class="btn btn-danger ml-2 remove-container">Remove</button>
          </div>`
    );
  });

  $("#add-deny-rule").click(function () {
    $("#deny-rules").append(
      `<div class="form-inline mb-2">
              <input type="text" class="form-control mr-2" name="deny_from" placeholder="From">
              <input type="text" class="form-control" name="deny_to" placeholder="To">
              <button type="button" class="btn btn-danger ml-2 remove-deny-rule">Remove</button>
          </div>`
    );
  });

  $("#add-high-risk-path").click(function () {
    $("#high-risk-paths").append(
      `<div class="form-inline mb-2">
              <input type="text" class="form-control" name="high_risk_paths" placeholder="Path">
              <button type="button" class="btn btn-danger ml-2 remove-high-risk-path">Remove</button>
          </div>`
    );
  });

  $("#add-syscall-id").click(function () {
    $("#syscall-blacklist").append(
      `<div class="form-inline mb-2">
              <input type="text" class="form-control" name="syscall_id" placeholder="Syscall ID">
              <button type="button" class="btn btn-danger ml-2 remove-syscall-id">Remove</button>
          </div>`
    );
  });

  $(document).on('click', '.remove-container', function () {
    $(this).parent().remove();
  });

  $(document).on('click', '.remove-deny-rule', function () {
    $(this).parent().remove();
  });

  $(document).on('click', '.remove-high-risk-path', function () {
    $(this).parent().remove();
  });

  $(document).on('click', '.remove-syscall-id', function () {
    $(this).parent().remove();
  });
});
