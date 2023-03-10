import React, { memo } from "react";
import { useDispatch } from "react-redux";

import {Button} from "tabler-react";
import {
  acceptOrganizationInvite as acceptOrganizationInviteAction,
  declineOrganizationInvite as declineOrganizationInviteAction
} from "../../actions/organizations";

const AcceptOrganizationInviteButton = ({ organizationInvite, organizationName }) => {
  const dispatch = useDispatch();
  const acceptOrganizationInvite = (organizationInviteID, organizationName) =>
    dispatch(acceptOrganizationInviteAction(organizationInviteID, organizationName));
  const declineOrganizationInvite = (organizationInviteID, organizationName) =>
    dispatch(declineOrganizationInviteAction(organizationInviteID, organizationName));

  const handleAcceptOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await acceptOrganizationInvite(organizationInvite.id, organizationName);
  };
  const handleDeclineOrganizationInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await declineOrganizationInvite(organizationInvite.id, organizationName);
  };

  return (
    <>
      <Button.List>
        <Button color="success" className="mx-lg-auto" onClick={handleAcceptOrganizationInvite}>✔</Button>
        <Button color="red" className="ml-auto" onClick={handleDeclineOrganizationInvite}>✖</Button>
      </Button.List>
    </>
  );
};

export default memo(AcceptOrganizationInviteButton);
