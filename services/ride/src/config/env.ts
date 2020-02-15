interface Environment {
  estimateExpirationTime: number;
  port: number;
  pricingAddr: string;
  userAddr: string;
}

const env: Environment = {
  estimateExpirationTime: 20,
  port: 3001,
  pricingAddr: "pricing-api:50051",
  userAddr: "user-api:50051"
};

export default env;
