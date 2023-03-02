import React, { memo } from "react";
import { useDispatch } from "react-redux";

import {Button} from "tabler-react";
import {
  acceptTeamInvite as acceptTeamInviteAction,
  declineTeamInvite as declineTeamInviteAction
} from "../../actions/seasons";

const AcceptTeamInviteButton = ({ teamInvite, seasonName, organizationName }) => {
  const dispatch = useDispatch();
  const acceptTeamInvite = (teamInviteID, seasonName, organizationName) =>
    dispatch(acceptTeamInviteAction(teamInviteID, seasonName, organizationName));
  const declineTeamInvite = (teamInviteID, seasonName, organizationName) =>
    dispatch(declineTeamInviteAction(teamInviteID, seasonName, organizationName));

  const handleAcceptTeamInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await acceptTeamInvite(teamInvite.id, seasonName,organizationName);
  };
  const handleDeclineTeamInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await declineTeamInvite(teamInvite.id, seasonName,organizationName);
  };

  return (
    <>
      <Button.List>
        <Button color="success" className="mx-lg-auto" onClick={handleAcceptTeamInvite}>✔</Button>
        <Button color="red" className="ml-auto" onClick={handleDeclineTeamInvite}>✖</Button>
      </Button.List>
    </>
  );
};

export default memo(AcceptTeamInviteButton);
