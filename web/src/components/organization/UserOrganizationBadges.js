import React from "react";
import {Avatar} from "tabler-react";
import {Link} from "gatsby";

const styles = {
  margin: "0.25rem",
}

export const UserOrganizationBadges = ({organizations}) => {
  return !organizations ? (
    <></>
    ) : (
      <Avatar.List>
        {organizations.map(organization => (
          <Link key={organization.id} to={"/organization/" + organization.id}>
            <Avatar size="md" imageURL={`${organization.gravatar_url}?d=identicon`} css={styles}/>
          </Link>
     ))}
      </Avatar.List>
  );
}

export default UserOrganizationBadges;
