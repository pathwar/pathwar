import React, { useEffect, useState } from "react";
import { Page, Grid, Avatar, Dimmer, ProgressCard } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import {
  fetchOrganizationDetail as fetchOrganizationDetailAction,
} from "../actions/organizations";
import {CLEAN_ORGANIZATION_DETAILS} from "../constants/actionTypes";

const ChallengeDetailsPage = props => {
  const dispatch = useDispatch();

  const organization = useSelector(state => state.organizations.organizationInDetail);

  const fetchOrganizationDetail = challengeID =>
    dispatch(fetchOrganizationDetailAction(challengeID));

  useEffect(() => {
    const { uri, organizationID: organizationIDFromProps } = props;
    const challengeIDFromURI = uri && uri.split("/")[2];
    const organizationID = challengeIDFromURI || organizationIDFromProps;

    fetchOrganizationDetail(organizationID);

    return () => dispatch({ type: CLEAN_ORGANIZATION_DETAILS });
   }, []);

  if (!organization) {
    return <Dimmer active loader />;
  }

  return (
    <Page.Content title={"TEST"}>
      <span>{organization.name}</span>
    </Page.Content>
  );
};

export default React.memo(ChallengeDetailsPage);
