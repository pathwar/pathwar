import React, { useState } from "react";
import { useDispatch } from "react-redux";
import { Helmet } from "react-helmet";
import { Page, Grid, Button } from "tabler-react";
import { useIntl, FormattedMessage } from "react-intl";

import siteMetaData from "../constants/metadata";
import { deleteAccount as deleteAccountAction } from "../actions/userSession";

const SettingsPage = () => {
  const intl = useIntl();
  const dispatch = useDispatch();
  const [isFetching, setFetching] = useState(false);
  const { title, description } = siteMetaData;

  const deleteAccountDispatch = reason => dispatch(deleteAccountAction(reason));

  const deleteAccount = async reason => {
    setFetching(true);
    deleteAccountDispatch(reason).then(response => {
      setFetching(false);
      return response;
    });
  };

  const pageTitleIntl = intl.formatMessage({ id: "SettingsPage.settings" });

  return (
    <>
      <Helmet>
        <title>
          {title} - {pageTitleIntl}
        </title>
        <meta name="description" content={description} />
      </Helmet>
      <Page.Content title={pageTitleIntl}>
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={12} lg={6}>
            <Button.List>
              <Button
                onClick={() => deleteAccount("integration test")}
                loading={isFetching}
                color="primary"
              >
                <FormattedMessage id="SettingsPage.deleteAccount" />
              </Button>
            </Button.List>
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    </>
  );
};

export default SettingsPage;
