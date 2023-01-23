def instance_placement(reason, request, candidate_members):
    # Example of logging info, this will appear in LXD's log.
    log_info("instance placement started: ", reason, request)

    # Example of applying logic based on the instance request.
    if request["name"] == "foo":
        # Example of logging an error, this will appear in LXD's log.
        log_error("Invalid name supplied: ", request["name"])

        return "Invalid name" # Return an error to reject instance placement.

    # Place the instance on the first candidate server provided.
    set_target(candidate_members[0]["server_name"])

    return # Return empty to allow instance placement to proceed.
