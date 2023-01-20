import React, { useEffect } from "react";
import { Page, Grid, Avatar, Dimmer } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import {
  fetchOrganizationDetail as fetchOrganizationDetailAction,
} from "../actions/organizations";
import {CLEAN_ORGANIZATION_DETAILS} from "../constants/actionTypes";
import ShadowBox from "../components/ShadowBox";
import {FormattedMessage} from "react-intl";
import moment from "moment/moment";
import UserOrganizationsList from "../components/organization/UserOrganizationsList";
import TeamsOnOrganizationList from "../components/organization/AllTeamsOnOrganization";

const OrganizationDetailsPage = props => {
  const dispatch = useDispatch();

  const organization = useSelector(state => state.organizations.organizationInDetail);

  const fetchOrganizationDetail = organizationID =>
    dispatch(fetchOrganizationDetailAction(organizationID));

  useEffect(() => {
    const { uri, organizationID: organizationIDFromProps } = props;
    const organizationIDFromURI = uri && uri.split("/")[2];
    const organizationID = organizationIDFromURI || organizationIDFromProps;

    fetchOrganizationDetail(organizationID);

    return () => dispatch({ type: CLEAN_ORGANIZATION_DETAILS });
   }, []);

  if (!organization) {
    return <Dimmer active loader />;
  }

  return (
    <Page.Content title={organization.name}>
      <Grid.Row>
        <Grid.Col width={12} lg={5}>
          <ShadowBox>
            <div
              css={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
              }}
            >
              <h2 className="mb-4 mt-2" style={{fontSize: '1.60rem'}}>{organization.name}</h2>
              {organization.gravatar_url ? (
                <Avatar size="xxl" imageURL={`${organization.gravatar_url}?d=identicon`} />
              ) : (
                <Avatar size="xxl" icon="users" />
              )}
              <h3 className="mb-0 mt-4">
                <FormattedMessage id="HomePage.createdAt" />
              </h3>
              <p >{moment(organization.created_at).format("ll")}</p>
            </div>
          </ShadowBox>
        </Grid.Col>
        <Grid.Col xs={12} sm={12} md={6}>
          <TeamsOnOrganizationList
            teams={organization.teams}
          />
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(OrganizationDetailsPage);
