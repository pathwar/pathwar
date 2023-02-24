import React, { useEffect } from "react";
import { Page, Grid, Avatar, Dimmer } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import {
  fetchOrganizationDetail as fetchOrganizationDetailAction,
} from "../../actions/organizations";
import {CLEAN_ORGANIZATION_DETAILS} from "../../constants/actionTypes";
import ShadowBox from "../../components/ShadowBox";
import {FormattedMessage, useIntl} from "react-intl";
import moment from "moment/moment";
import TeamsOnOrganizationList from "../../components/organization/AllTeamsOnOrganization";
import {css} from "@emotion/core";
import OrganizationSubMenu from "../../components/organization/OrganizationSubMenu";
import InviteOrganizationButton from "../../components/organization/InviteOrganizationButton";

const wrapper = css`

.link {
  display: block;
  text-decoration: none;
  padding: 1rem;
  color: #919aa3;
  font-size: 1.2rem;

&:hover {
    opacity: 0.8;
  }
}

`
const OrganizationDetailsPage = props => {
  const intl = useIntl();
  const pageTitleIntl = intl.formatMessage({ id: "OrganizationsPage.title" });

  const dispatch = useDispatch();
  const organization = useSelector(state => state.organizations.organizationInDetail);

  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
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

  const adminActions = activeUserSession && organization.members.filter(member => member.role === "Owner" && member.user.id === activeUserSession.user.id).length > 0;

  return (
    <Page.Content title={pageTitleIntl} css={wrapper}>
      <Grid.Row css={{
        "margin-bottom": "15px",
      }}>
        <OrganizationSubMenu organization={organization} />
      </Grid.Row>
      <Grid.Row>
        <Grid.Col width={12} lg={5} >
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
              <h3 className="mb-0 mt-6">
                <FormattedMessage id="OrganizationDetailsSubmenu.members" />
              </h3>
              <p >{organization.members.length}</p>
              <h3 className="mb-0 mt-4">
                <FormattedMessage id="HomePage.createdAt" />
              </h3>
              <p >{moment(organization.created_at).format("ll")}</p>
              {adminActions ?
                <>
                <InviteOrganizationButton organizationID={organization.id} organizationName={organization.name}/>
                <InviteOrganizationButton organizationID={organization.id} organizationName={organization.name}/>
                </>: null
              }
            </div>
          </ShadowBox>
        </Grid.Col>
        <Grid.Col xs={12} sm={12} md={6}>
          <TeamsOnOrganizationList
            teams={organization.teams}/>
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(OrganizationDetailsPage);
