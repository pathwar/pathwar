import React, { memo } from "react";
import { useDispatch } from "react-redux";

import {Button} from "tabler-react";
import {
  acceptOrganizationInvite as acceptOrganizationInviteAction,
  rejectOrganizationInvite as rejectOrganizationInviteAction
} from "../../actions/organizations";

const AcceptTeamInviteButton = ({ teamInvite, organizationName }) => {
  const dispatch = useDispatch();
  const acceptOrganizationInvite = (organizationInviteID, organizationName) =>
    dispatch(acceptOrganizationInviteAction(organizationInviteID, organizationName));
  const rejectOrganizationInvite = (organizationInviteID, organizationName) =>
    dispatch(rejectOrganizationInviteAction(organizationInviteID, organizationName));

  const handleAcceptOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await acceptOrganizationInvite(teamInvite.id, organizationName);
  };
  const handleRejectOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await rejectOrganizationInvite(teamInvite.id, organizationName);
  };

  return (
    <>
      <Button.List>
        <Button color="success" className="mx-lg-auto" onClick={handleAcceptOrganizationInvite}>✔</Button>
        <Button color="red" className="ml-auto" onClick={handleRejectOrganizationInvite}>✖</Button>
      </Button.List>
    </>
  );
};

export default memo(AcceptOrganizationInviteButton);
