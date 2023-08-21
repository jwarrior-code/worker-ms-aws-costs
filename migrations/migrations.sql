CREATE TABLE aws_storage_products (
                                      sku varchar(255) NOT NULL,
                                      product_family varchar(255) NULL,
                                      servicecode varchar(255) NULL,
                                      location varchar(255) NULL,
                                      volume_type varchar(255) NULL,
                                      price_per_gb_month numeric(10, 4) NULL,
                                      term_type varchar(255) NULL,
                                      billing_item varchar(255) NULL,
                                      price_per_hour numeric(10, 4) NULL,
                                      CONSTRAINT ebs_products_pkey PRIMARY KEY (sku)
);

CREATE TABLE aws_ec2_products (
                                  sku varchar(255) NOT NULL,
                                  product_family varchar(255) NULL,
                                  servicecode varchar(255) NULL,
                                  location varchar(255) NULL,
                                  volume_type varchar(255) NULL,
                                  instance_type varchar(255) NULL,
                                  vcpu int4 NULL,
                                  memory varchar(255) NULL,
                                  network_performance varchar(255) NULL,
                                  term_type varchar(255) NULL,
                                  price_per_hour numeric(10, 6) NULL,
                                  billing_item varchar(255) NULL,
                                  storage varchar(255) NULL,
                                  CONSTRAINT ec2_products_pkey PRIMARY KEY (sku)
);

CREATE TABLE aws_rds_products (
                                  sku varchar(255) NOT NULL,
                                  product_family varchar(255) NULL,
                                  servicecode varchar(255) NULL,
                                  location varchar(255) NULL,
                                  database_engine varchar(255) NULL,
                                  deployment_option varchar(255) NULL,
                                  license_model varchar(255) NULL,
                                  vcpu int4 NULL,
                                  memory varchar(255) NULL,
                                  instance_class varchar(255) NULL,
                                  term_type varchar(255) NULL,
                                  price_per_hour numeric(10, 6) NULL,
                                  billing_item varchar(255) NULL,
                                  CONSTRAINT rds_products_pkey PRIMARY KEY (sku)
);
