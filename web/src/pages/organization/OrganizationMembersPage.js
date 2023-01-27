import React, { useEffect } from "react";
import { Page, Grid, Avatar, Dimmer } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import {
  fetchOrganizationDetail as fetchOrganizationDetailAction,
} from "../../actions/organizations";
import {CLEAN_ORGANIZATION_DETAILS} from "../../constants/actionTypes";
import {useIntl} from "react-intl";
import {css} from "@emotion/core";
import OrganizationSubMenu from "../../components/organization/OrganizationSubMenu";
import MembersOnOrganizationList from "../../components/organization/AllMembersOnOrganization";

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

const OrganizationMembersPage = props => {
  const intl = useIntl();
  const pageTitleIntl = intl.formatMessage({ id: "OrganizationsPage.title" });

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
    <Page.Content title={pageTitleIntl} css={wrapper}>
      <Grid.Row css={{
        "margin-bottom": "15px",
      }}>
        <OrganizationSubMenu organization={organization} />
      </Grid.Row>
      <Grid.Row>
        <Grid.Col xs={12} sm={12} md={12}>
          <MembersOnOrganizationList
            members={organization.members}
          />
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(OrganizationMembersPage);
