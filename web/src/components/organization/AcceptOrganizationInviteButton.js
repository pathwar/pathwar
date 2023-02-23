import React, { memo } from "react";
import { useDispatch } from "react-redux";

import {Button} from "tabler-react";
import {
  acceptOrganizationInvite as acceptOrganizationInviteAction,
  rejectOrganizationInvite as rejectOrganizationInviteAction
} from "../../actions/organizations";

const AcceptOrganizationInviteButton = ({ organizationInvite, ...rest }) => {
  const dispatch = useDispatch();
  const acceptOrganizationInvite = (organizationInviteID) =>
    dispatch(acceptOrganizationInviteAction(organizationInviteID));
  const rejectOrganizationInvite = (organizationInviteID) =>
    dispatch(rejectOrganizationInviteAction(organizationInviteID));

  const handleAcceptOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await acceptOrganizationInvite(organizationInvite.id, true);
  };
  const handleRejectOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await rejectOrganizationInvite(organizationInvite.id, false);
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
