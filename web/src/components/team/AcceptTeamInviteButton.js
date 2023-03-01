import React, { memo } from "react";
import { useDispatch } from "react-redux";

import {Button} from "tabler-react";
import {
  acceptTeamInvite as acceptTeamInviteAction,
  rejectTeamInvite as rejectTeamInviteAction
} from "../../actions/seasons";

const AcceptTeamInviteButton = ({ teamInvite, seasonName, organizationName }) => {
  const dispatch = useDispatch();
  const acceptTeamInvite = (teamInviteID, seasonName, organizationName) =>
    dispatch(acceptTeamInviteAction(teamInviteID, seasonName, organizationName));
  const rejectTeamInvite = (teamInviteID, seasonName, organizationName) =>
    dispatch(rejectTeamInviteAction(teamInviteID, seasonName, organizationName));

  const handleAcceptTeamInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await acceptTeamInvite(teamInvite.id, seasonName,organizationName);
  };
  const handleRejectTeamInvite = async event => {
    event.preventDefault();
    event.stopPropagation();
    await rejectTeamInvite(teamInvite.id, seasonName,organizationName);
  };

  return (
    <>
      <Button.List>
        <Button color="success" className="mx-lg-auto" onClick={handleAcceptTeamInvite}>✔</Button>
        <Button color="red" className="ml-auto" onClick={handleRejectTeamInvite}>✖</Button>
      </Button.List>
    </>
  );
};

export default memo(AcceptTeamInviteButton);
