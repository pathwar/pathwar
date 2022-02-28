import React from "react";
import { Page, Grid } from "tabler-react";
import { FormattedMessage, useIntl } from "react-intl";

const ComingSoonPage = () => {
  const intl = useIntl();

  const pageTitleIntl = intl.formatMessage({ id: "ComingSoonPage.title" });

  return (
    <Page.Content title={pageTitleIntl}>
      <Grid.Row>
        <Grid.Col xs={12} sm={12} lg={6}>
          <h2>
            <FormattedMessage id="ComingSoonPage.text" />
            ...
          </h2>
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(ComingSoonPage);
