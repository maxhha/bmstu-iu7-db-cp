import { GraphQLClient } from "graphql-request"
import * as Dom from "graphql-request/dist/types.dom"
import { print } from "graphql"
import gql from "graphql-tag"
export type Maybe<T> = T | null
export type InputMaybe<T> = Maybe<T>
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K]
}
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>
}
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>
}
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string
  String: string
  Boolean: boolean
  Int: number
  Float: number
  Cursor: any
  DateTime: any
  Map: any
}

/** Nominal account */
export type Account = {
  /** Bank in which the account was created */
  bank: Bank
  id: Scalars["ID"]
  /** All transactions in which the account is involved */
  transactions: TransactionsConnection
}

/** Nominal account */
export type AccountTransactionsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type AccountsConnection = {
  __typename?: "AccountsConnection"
  edges: Array<AccountsConnectionEdge>
  pageInfo: PageInfo
}

export type AccountsConnectionEdge = {
  __typename?: "AccountsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: Account
}

export type ApproveUserFormInput = {
  userFormId: Scalars["ID"]
}

export type Auction = {
  __typename?: "Auction"
  /** Product buyer */
  buyer?: Maybe<User>
  /** Auctions currency */
  currency: CurrencyEnum
  /** Real time of auction end */
  finishedAt?: Maybe<Scalars["DateTime"]>
  id: Scalars["String"]
  /** First offer must have equal or greater amount */
  minAmount?: Maybe<Scalars["Float"]>
  /** Product for selling */
  product: Product
  /** Planned time for auction end */
  scheduledFinishAt?: Maybe<Scalars["DateTime"]>
  /** Planned time for auction start */
  scheduledStartAt?: Maybe<Scalars["DateTime"]>
  /** Product seller, auction creator */
  seller: User
  /** Real auction start time */
  startedAt?: Maybe<Scalars["DateTime"]>
  state: AuctionState
}

export type AuctionInput = {
  auctionId: Scalars["ID"]
}

export type AuctionResult = {
  __typename?: "AuctionResult"
  auction: Auction
}

export enum AuctionState {
  Created = "CREATED",
  Failed = "FAILED",
  Finished = "FINISHED",
  Started = "STARTED",
  Succeeded = "SUCCEEDED",
}

export type AuctionsConnection = {
  __typename?: "AuctionsConnection"
  edges: Array<AuctionsConnectionEdge>
  pageInfo: PageInfo
}

export type AuctionsConnectionEdge = {
  __typename?: "AuctionsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: Auction
}

export type AuctionsFilter = {
  IDs?: InputMaybe<Array<Scalars["ID"]>>
  buyerIDs?: InputMaybe<Array<Scalars["String"]>>
  productIDs?: InputMaybe<Array<Scalars["String"]>>
  sellerIDs?: InputMaybe<Array<Scalars["String"]>>
  states?: InputMaybe<Array<AuctionState>>
}

/** Bank that is cooperated with platform */
export type Bank = {
  __typename?: "Bank"
  /** Special account of that bank */
  account: BankAccount
  id: Scalars["ID"]
  /** Name of bank */
  name: Scalars["String"]
}

/** Special account for banks. Amount on this account is always nonpositve */
export type BankAccount = Account & {
  __typename?: "BankAccount"
  /** Owner of account. Each bank have one special account */
  bank: Bank
  id: Scalars["ID"]
  /** All transactions in which the account is involved */
  transactions: TransactionsConnection
}

/** Special account for banks. Amount on this account is always nonpositve */
export type BankAccountTransactionsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type CreateOfferInput = {
  amount: Scalars["Float"]
  productId: Scalars["ID"]
}

export type CreateOfferResult = {
  __typename?: "CreateOfferResult"
  offer: Offer
}

export enum CurrencyEnum {
  Eur = "EUR",
  Rub = "RUB",
  Usd = "USD",
}

export type DateTimeRange = {
  from?: InputMaybe<Scalars["DateTime"]>
  to?: InputMaybe<Scalars["DateTime"]>
}

export type DeclineProductInput = {
  declainReason?: InputMaybe<Scalars["String"]>
  productId: Scalars["ID"]
}

export type DeclineUserFormInput = {
  declainReason?: InputMaybe<Scalars["String"]>
  userFormId: Scalars["ID"]
}

export type LoginInput = {
  password: Scalars["String"]
  username: Scalars["String"]
}

/** Мoney in a specific currency */
export type Money = {
  __typename?: "Money"
  amount: Scalars["Float"]
  currency: CurrencyEnum
}

/** Input money in a specific currency */
export type MoneyInput = {
  amount: Scalars["Float"]
  currency: CurrencyEnum
}

export type Mutation = {
  __typename?: "Mutation"
  /** Set product state to moderating */
  approveModerateProduct: ProductResult
  /** Set user form state to moderate */
  approveModerateUserForm: UserResult
  /** Approve product */
  approveProduct: ProductResult
  /** First input of users email */
  approveSetUserEmail: UserResult
  /** First input of users */
  approveSetUserPhone: UserResult
  /** Approve user form */
  approveUserForm: UserFormResult
  /** Create auction for given product */
  createAuction: AuctionResult
  createOffer: CreateOfferResult
  /** Creates product with creator of current viewer */
  createProduct: ProductResult
  /** Declain product */
  declainProduct: ProductResult
  /** Decline user form */
  declineUserForm: UserFormResult
  /** User login */
  login: TokenResult
  offerProduct: OfferProductResult
  /** Registrates empty user */
  register: TokenResult
  removeOffer: RemoveOfferResult
  /** Request token to send product for modetaion */
  requestModerateProduct: Scalars["Boolean"]
  /** Send token for user form moderation */
  requestModerateUserForm: Scalars["Boolean"]
  /** Request set user email */
  requestSetUserEmail: Scalars["Boolean"]
  /** Request set user email */
  requestSetUserPhone: Scalars["Boolean"]
  sellProduct: SellProductResult
  /** Starts auction manually */
  startAuction: AuctionResult
  takeOffProduct: TakeOffProductResult
  /** Update auction */
  updateAuction: AuctionResult
  /** Update product info */
  updateProduct: ProductResult
  /** Update user draft form fields */
  updateUserDraftForm: UserResult
  /** Update user password using old password */
  updateUserPassword: UserResult
}

export type MutationApproveModerateProductArgs = {
  input: TokenInput
}

export type MutationApproveModerateUserFormArgs = {
  input: TokenInput
}

export type MutationApproveProductArgs = {
  input: ProductInput
}

export type MutationApproveSetUserEmailArgs = {
  input: TokenInput
}

export type MutationApproveSetUserPhoneArgs = {
  input: TokenInput
}

export type MutationApproveUserFormArgs = {
  input: ApproveUserFormInput
}

export type MutationCreateAuctionArgs = {
  input: ProductInput
}

export type MutationCreateOfferArgs = {
  input: CreateOfferInput
}

export type MutationDeclainProductArgs = {
  input: DeclineProductInput
}

export type MutationDeclineUserFormArgs = {
  input: DeclineUserFormInput
}

export type MutationLoginArgs = {
  input: LoginInput
}

export type MutationOfferProductArgs = {
  input: ProductInput
}

export type MutationRemoveOfferArgs = {
  input: RemoveOfferInput
}

export type MutationRequestModerateProductArgs = {
  input: ProductInput
}

export type MutationRequestSetUserEmailArgs = {
  input: RequestSetUserEmailInput
}

export type MutationRequestSetUserPhoneArgs = {
  input: RequestSetUserPhoneInput
}

export type MutationSellProductArgs = {
  input: ProductInput
}

export type MutationStartAuctionArgs = {
  input: AuctionInput
}

export type MutationTakeOffProductArgs = {
  input: ProductInput
}

export type MutationUpdateAuctionArgs = {
  input: UpdateAuctionInput
}

export type MutationUpdateProductArgs = {
  input: UpdateProductInput
}

export type MutationUpdateUserDraftFormArgs = {
  input: UpdateUserDraftFormInput
}

export type MutationUpdateUserPasswordArgs = {
  input: UpdateUserPasswordInput
}

export type Offer = {
  __typename?: "Offer"
  /** Offer creation time */
  createdAt: Scalars["DateTime"]
  /** If set to true, the offer will be removed after the product is sold */
  deleteOnSell: Scalars["Boolean"]
  /** Reason of fail for *_FAILED states */
  failReason?: Maybe<Scalars["String"]>
  id: Scalars["ID"]
  /** Total moneys offered */
  moneys: Array<Money>
  /** Product for which this offer was created */
  product: Product
  /** Current offer state */
  state: OfferStateEnum
  /** Transactions of this offer */
  transactions: Array<Transaction>
  /** User created this offer */
  user: User
}

export type OfferProductResult = {
  __typename?: "OfferProductResult"
  product: Product
}

export enum OfferStateEnum {
  Cancelled = "CANCELLED",
  Created = "CREATED",
  MoneyReturned = "MONEY_RETURNED",
  ReturningMoney = "RETURNING_MONEY",
  ReturnMoneyFailed = "RETURN_MONEY_FAILED",
  Succeeded = "SUCCEEDED",
  TransferringMoney = "TRANSFERRING_MONEY",
  TransferringProduct = "TRANSFERRING_PRODUCT",
  TransferMoneyFailed = "TRANSFER_MONEY_FAILED",
  TransferProductFailed = "TRANSFER_PRODUCT_FAILED",
}

export type OffersConnection = {
  __typename?: "OffersConnection"
  edges: Array<OffersConnectionEdge>
  pageInfo: PageInfo
}

export type OffersConnectionEdge = {
  __typename?: "OffersConnectionEdge"
  cursor: Scalars["Cursor"]
  node: Offer
}

export type PageInfo = {
  __typename?: "PageInfo"
  endCursor?: Maybe<Scalars["Cursor"]>
  hasNextPage: Scalars["Boolean"]
  hasPreviousPage: Scalars["Boolean"]
  startCursor?: Maybe<Scalars["Cursor"]>
}

export type Product = {
  __typename?: "Product"
  /** Creator of product */
  creator: User
  /** Declain reason */
  declainReason?: Maybe<Scalars["String"]>
  /** Product description */
  description: Scalars["String"]
  id: Scalars["ID"]
  /** Product images */
  images: Array<ProductImage>
  /** Offers for this product */
  offers: OffersConnection
  /** Current owner of product */
  owner: User
  /** Current state of product */
  state: ProductState
  /** Title of product */
  title: Scalars["String"]
  /** The greatest offer */
  topOffer?: Maybe<Offer>
}

export type ProductOffersArgs = {
  after?: InputMaybe<Scalars["String"]>
  first?: InputMaybe<Scalars["Int"]>
}

/** Product image */
export type ProductImage = {
  __typename?: "ProductImage"
  filename: Scalars["String"]
  id: Scalars["ID"]
  path: Scalars["String"]
}

export type ProductInput = {
  productId: Scalars["ID"]
}

export type ProductResult = {
  __typename?: "ProductResult"
  product: Product
}

export enum ProductState {
  Approved = "APPROVED",
  Created = "CREATED",
  Declained = "DECLAINED",
  Moderating = "MODERATING",
}

export type ProductsConnection = {
  __typename?: "ProductsConnection"
  edges: Array<ProductsConnectionEdge>
  pageInfo: PageInfo
}

export type ProductsConnectionEdge = {
  __typename?: "ProductsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: Product
}

export type ProductsFilter = {
  ownerIDs?: InputMaybe<Array<Scalars["String"]>>
}

export type Query = {
  __typename?: "Query"
  /** All auctions */
  auctions: AuctionsConnection
  marketProducts: ProductsConnection
  /** List all products */
  products: ProductsConnection
  /** List of all user forms */
  userForms: UserFormsConnection
  /** List of all users */
  users?: Maybe<UsersConnection>
  /** Authorized user */
  viewer?: Maybe<User>
}

export type QueryAuctionsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  filter?: InputMaybe<AuctionsFilter>
  first?: InputMaybe<Scalars["Int"]>
}

export type QueryMarketProductsArgs = {
  after?: InputMaybe<Scalars["String"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type QueryProductsArgs = {
  after?: InputMaybe<Scalars["String"]>
  filter?: InputMaybe<ProductsFilter>
  first?: InputMaybe<Scalars["Int"]>
}

export type QueryUserFormsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  filter?: InputMaybe<UserFormsFilter>
  first?: InputMaybe<Scalars["Int"]>
}

export type QueryUsersArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  filter?: InputMaybe<UsersFilter>
  first?: InputMaybe<Scalars["Int"]>
}

export type RemoveOfferInput = {
  offerId: Scalars["ID"]
}

export type RemoveOfferResult = {
  __typename?: "RemoveOfferResult"
  status: Scalars["String"]
}

export type RequestSetUserEmailInput = {
  email: Scalars["String"]
}

export type RequestSetUserPhoneInput = {
  phone: Scalars["String"]
}

export enum RoleType {
  Admin = "ADMIN",
  Manager = "MANAGER",
}

export type SellProductResult = {
  __typename?: "SellProductResult"
  product: Product
}

export type Subscription = {
  __typename?: "Subscription"
  productOffered?: Maybe<Product>
}

export type TakeOffProductResult = {
  __typename?: "TakeOffProductResult"
  product: Product
}

/** Used for actions activation */
export type TokenInput = {
  token: Scalars["String"]
}

/** Used for login and registration */
export type TokenResult = {
  __typename?: "TokenResult"
  token: Scalars["String"]
}

export type Transaction = {
  __typename?: "Transaction"
  /** From account */
  accountFrom: Account
  /** To account */
  accountTo: Account
  /** Transaction amount */
  amount: Scalars["Float"]
  /** Transaction currency */
  currency: CurrencyEnum
  /** Time of apply this transaction */
  date?: Maybe<Scalars["DateTime"]>
  /** Error message for state = ERROR or FAILED */
  error?: Maybe<Scalars["String"]>
  id: Scalars["ID"]
  /** Offer for type = BUY */
  offer?: Maybe<Offer>
  /** Current state */
  state: TransactionStateEnum
  /** Transaction type */
  type: TransactionTypeEnum
}

export enum TransactionStateEnum {
  Cancelled = "CANCELLED",
  Created = "CREATED",
  Error = "ERROR",
  Failed = "FAILED",
  Processing = "PROCESSING",
  Succeeded = "SUCCEEDED",
}

export enum TransactionTypeEnum {
  Buy = "BUY",
  Deposit = "DEPOSIT",
  Fee = "FEE",
  Withdrawal = "WITHDRAWAL",
}

export type TransactionsConnection = {
  __typename?: "TransactionsConnection"
  edges: Array<TransactionsConnectionEdge>
  pageInfo: PageInfo
}

export type TransactionsConnectionEdge = {
  __typename?: "TransactionsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: Transaction
}

export type UpdateAuctionInput = {
  auctionId: Scalars["ID"]
  currency: CurrencyEnum
  minAmount?: InputMaybe<Scalars["Float"]>
  scheduledFinishAt?: InputMaybe<Scalars["DateTime"]>
  scheduledStartAt?: InputMaybe<Scalars["DateTime"]>
}

export type UpdateProductInput = {
  description: Scalars["String"]
  productId: Scalars["ID"]
  title: Scalars["String"]
}

export type UpdateUserDraftFormInput = {
  currency?: InputMaybe<CurrencyEnum>
  name?: InputMaybe<Scalars["String"]>
}

export type UpdateUserPasswordInput = {
  oldPassword?: InputMaybe<Scalars["String"]>
  password: Scalars["String"]
}

export type User = {
  __typename?: "User"
  /** User accounts */
  accounts: UserAccountsConnection
  /** Auctions which user created */
  auctions: AuctionsConnection
  /** Available moneys */
  available: Array<Money>
  /** Money that is blocked in some offers */
  blocked: Array<Money>
  /** End date of blocking this user */
  blockedUntil?: Maybe<Scalars["DateTime"]>
  /** User new personal information */
  draftForm?: Maybe<UserForm>
  /** User current personal information */
  form?: Maybe<UserFormFilled>
  /** User history of personal information (only for managers) */
  formHistory: UserFormsConnection
  id: Scalars["ID"]
  /** User offers */
  offers: OffersConnection
  /** User products in which he is owner */
  products: ProductsConnection
}

export type UserAccountsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type UserAuctionsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type UserFormHistoryArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  filter?: InputMaybe<UserFormHistoryFilter>
  first?: InputMaybe<Scalars["Int"]>
}

export type UserOffersArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type UserProductsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

/** Nominal account that was created for client */
export type UserAccount = Account & {
  __typename?: "UserAccount"
  /** Bank in which the account was created */
  bank: Bank
  id: Scalars["ID"]
  /** All transactions in which the account is involved */
  transactions: TransactionsConnection
  /** Owner of account */
  user: User
}

/** Nominal account that was created for client */
export type UserAccountTransactionsArgs = {
  after?: InputMaybe<Scalars["Cursor"]>
  first?: InputMaybe<Scalars["Int"]>
}

export type UserAccountsConnection = {
  __typename?: "UserAccountsConnection"
  edges: Array<UserAccountsConnectionEdge>
  pageInfo: PageInfo
}

/** Connection with UserAccount only */
export type UserAccountsConnectionEdge = {
  __typename?: "UserAccountsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: UserAccount
}

/** User personal information */
export type UserForm = {
  __typename?: "UserForm"
  /** User default currency */
  currency?: Maybe<CurrencyEnum>
  /** User email */
  email?: Maybe<Scalars["String"]>
  id: Scalars["ID"]
  /** User name */
  name?: Maybe<Scalars["String"]>
  /** User phone */
  phone?: Maybe<Scalars["String"]>
  /** User form state */
  state: UserFormState
  /** UserForm owner */
  user: User
}

/** UserFrom with all required fields filled in */
export type UserFormFilled = {
  __typename?: "UserFormFilled"
  /** User default currency */
  currency: CurrencyEnum
  /** User email */
  email: Scalars["String"]
  /** User name */
  name: Scalars["String"]
  /** User phone */
  phone: Scalars["String"]
}

export type UserFormHistoryFilter = {
  id?: InputMaybe<Array<Scalars["ID"]>>
  state?: InputMaybe<Array<UserFormState>>
}

export type UserFormResult = {
  __typename?: "UserFormResult"
  userForm: UserForm
}

export enum UserFormState {
  Approved = "APPROVED",
  Created = "CREATED",
  Declained = "DECLAINED",
  Moderating = "MODERATING",
}

export type UserFormsConnection = {
  __typename?: "UserFormsConnection"
  edges: Array<UserFormsConnectionEdge>
  pageInfo: PageInfo
}

export type UserFormsConnectionEdge = {
  __typename?: "UserFormsConnectionEdge"
  cursor: Scalars["Cursor"]
  node: UserForm
}

export type UserFormsFilter = {
  id?: InputMaybe<Array<Scalars["ID"]>>
  state?: InputMaybe<Array<UserFormState>>
  userId?: InputMaybe<Array<Scalars["ID"]>>
}

export type UserResult = {
  __typename?: "UserResult"
  user: User
}

export type UsersConnection = {
  __typename?: "UsersConnection"
  edges: Array<UsersConnectionEdge>
  pageInfo: PageInfo
}

export type UsersConnectionEdge = {
  __typename?: "UsersConnectionEdge"
  cursor: Scalars["Cursor"]
  node: User
}

export type UsersFilter = {
  id?: InputMaybe<Array<Scalars["ID"]>>
}

export type RegisterMutationVariables = Exact<{ [key: string]: never }>

export type RegisterMutation = {
  __typename?: "Mutation"
  register: { __typename?: "TokenResult"; token: string }
}

export type ViewerQueryVariables = Exact<{ [key: string]: never }>

export type ViewerQuery = {
  __typename?: "Query"
  viewer?: { __typename?: "User"; id: string } | null
}

export type UpdateUserPasswordMutationVariables = Exact<{
  input: UpdateUserPasswordInput
}>

export type UpdateUserPasswordMutation = {
  __typename?: "Mutation"
  updateUserPassword: {
    __typename?: "UserResult"
    user: { __typename?: "User"; id: string }
  }
}

export type RequestSetUserEmailMutationVariables = Exact<{
  input: RequestSetUserEmailInput
}>

export type RequestSetUserEmailMutation = {
  __typename?: "Mutation"
  requestSetUserEmail: boolean
}

export type RequestSetUserPhoneMutationVariables = Exact<{
  input: RequestSetUserPhoneInput
}>

export type RequestSetUserPhoneMutation = {
  __typename?: "Mutation"
  requestSetUserPhone: boolean
}

export type ApproveSetUserEmailMutationVariables = Exact<{
  input: TokenInput
}>

export type ApproveSetUserEmailMutation = {
  __typename?: "Mutation"
  approveSetUserEmail: {
    __typename?: "UserResult"
    user: {
      __typename?: "User"
      draftForm?: { __typename?: "UserForm"; email?: string | null } | null
    }
  }
}

export type UpdateUserDraftFormMutationVariables = Exact<{
  input: UpdateUserDraftFormInput
}>

export type UpdateUserDraftFormMutation = {
  __typename?: "Mutation"
  updateUserDraftForm: {
    __typename?: "UserResult"
    user: {
      __typename?: "User"
      draftForm?: { __typename?: "UserForm"; name?: string | null } | null
    }
  }
}

export type ApproveSetUserPhoneMutationVariables = Exact<{
  input: TokenInput
}>

export type ApproveSetUserPhoneMutation = {
  __typename?: "Mutation"
  approveSetUserPhone: {
    __typename?: "UserResult"
    user: {
      __typename?: "User"
      draftForm?: { __typename?: "UserForm"; phone?: string | null } | null
    }
  }
}

export type LoginMutationVariables = Exact<{
  input: LoginInput
}>

export type LoginMutation = {
  __typename?: "Mutation"
  login: { __typename?: "TokenResult"; token: string }
}

export type RequestModerateUserFormMutationVariables = Exact<{
  [key: string]: never
}>

export type RequestModerateUserFormMutation = {
  __typename?: "Mutation"
  requestModerateUserForm: boolean
}

export type ApproveModerateUserFormMutationVariables = Exact<{
  input: TokenInput
}>

export type ApproveModerateUserFormMutation = {
  __typename?: "Mutation"
  approveModerateUserForm: {
    __typename?: "UserResult"
    user: {
      __typename?: "User"
      draftForm?: { __typename?: "UserForm"; state: UserFormState } | null
    }
  }
}

export const RegisterDocument = gql`
  mutation Register {
    register {
      token
    }
  }
`
export const ViewerDocument = gql`
  query Viewer {
    viewer {
      id
    }
  }
`
export const UpdateUserPasswordDocument = gql`
  mutation UpdateUserPassword($input: UpdateUserPasswordInput!) {
    updateUserPassword(input: $input) {
      user {
        id
      }
    }
  }
`
export const RequestSetUserEmailDocument = gql`
  mutation RequestSetUserEmail($input: RequestSetUserEmailInput!) {
    requestSetUserEmail(input: $input)
  }
`
export const RequestSetUserPhoneDocument = gql`
  mutation RequestSetUserPhone($input: RequestSetUserPhoneInput!) {
    requestSetUserPhone(input: $input)
  }
`
export const ApproveSetUserEmailDocument = gql`
  mutation ApproveSetUserEmail($input: TokenInput!) {
    approveSetUserEmail(input: $input) {
      user {
        draftForm {
          email
        }
      }
    }
  }
`
export const UpdateUserDraftFormDocument = gql`
  mutation UpdateUserDraftForm($input: UpdateUserDraftFormInput!) {
    updateUserDraftForm(input: $input) {
      user {
        draftForm {
          name
        }
      }
    }
  }
`
export const ApproveSetUserPhoneDocument = gql`
  mutation ApproveSetUserPhone($input: TokenInput!) {
    approveSetUserPhone(input: $input) {
      user {
        draftForm {
          phone
        }
      }
    }
  }
`
export const LoginDocument = gql`
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      token
    }
  }
`
export const RequestModerateUserFormDocument = gql`
  mutation RequestModerateUserForm {
    requestModerateUserForm
  }
`
export const ApproveModerateUserFormDocument = gql`
  mutation ApproveModerateUserForm($input: TokenInput!) {
    approveModerateUserForm(input: $input) {
      user {
        draftForm {
          state
        }
      }
    }
  }
`

export type SdkFunctionWrapper = <T>(
  action: (requestHeaders?: Record<string, string>) => Promise<T>,
  operationName: string,
  operationType?: string
) => Promise<T>

const defaultWrapper: SdkFunctionWrapper = (
  action,
  _operationName,
  _operationType
) => action()
export const RegisterDocumentString = print(RegisterDocument)
export const ViewerDocumentString = print(ViewerDocument)
export const UpdateUserPasswordDocumentString = print(UpdateUserPasswordDocument)
export const RequestSetUserEmailDocumentString = print(RequestSetUserEmailDocument)
export const RequestSetUserPhoneDocumentString = print(RequestSetUserPhoneDocument)
export const ApproveSetUserEmailDocumentString = print(ApproveSetUserEmailDocument)
export const UpdateUserDraftFormDocumentString = print(UpdateUserDraftFormDocument)
export const ApproveSetUserPhoneDocumentString = print(ApproveSetUserPhoneDocument)
export const LoginDocumentString = print(LoginDocument)
export const RequestModerateUserFormDocumentString = print(
  RequestModerateUserFormDocument
)
export const ApproveModerateUserFormDocumentString = print(
  ApproveModerateUserFormDocument
)
export function getSdk(
  client: GraphQLClient,
  withWrapper: SdkFunctionWrapper = defaultWrapper
) {
  return {
    Register(
      variables?: RegisterMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: RegisterMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<RegisterMutation>(
            RegisterDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "Register",
        "mutation"
      )
    },
    Viewer(
      variables?: ViewerQueryVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: ViewerQuery
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<ViewerQuery>(ViewerDocumentString, variables, {
            ...requestHeaders,
            ...wrappedRequestHeaders,
          }),
        "Viewer",
        "query"
      )
    },
    UpdateUserPassword(
      variables: UpdateUserPasswordMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: UpdateUserPasswordMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<UpdateUserPasswordMutation>(
            UpdateUserPasswordDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "UpdateUserPassword",
        "mutation"
      )
    },
    RequestSetUserEmail(
      variables: RequestSetUserEmailMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: RequestSetUserEmailMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<RequestSetUserEmailMutation>(
            RequestSetUserEmailDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "RequestSetUserEmail",
        "mutation"
      )
    },
    RequestSetUserPhone(
      variables: RequestSetUserPhoneMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: RequestSetUserPhoneMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<RequestSetUserPhoneMutation>(
            RequestSetUserPhoneDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "RequestSetUserPhone",
        "mutation"
      )
    },
    ApproveSetUserEmail(
      variables: ApproveSetUserEmailMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: ApproveSetUserEmailMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<ApproveSetUserEmailMutation>(
            ApproveSetUserEmailDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "ApproveSetUserEmail",
        "mutation"
      )
    },
    UpdateUserDraftForm(
      variables: UpdateUserDraftFormMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: UpdateUserDraftFormMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<UpdateUserDraftFormMutation>(
            UpdateUserDraftFormDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "UpdateUserDraftForm",
        "mutation"
      )
    },
    ApproveSetUserPhone(
      variables: ApproveSetUserPhoneMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: ApproveSetUserPhoneMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<ApproveSetUserPhoneMutation>(
            ApproveSetUserPhoneDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "ApproveSetUserPhone",
        "mutation"
      )
    },
    Login(
      variables: LoginMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: LoginMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<LoginMutation>(LoginDocumentString, variables, {
            ...requestHeaders,
            ...wrappedRequestHeaders,
          }),
        "Login",
        "mutation"
      )
    },
    RequestModerateUserForm(
      variables?: RequestModerateUserFormMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: RequestModerateUserFormMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<RequestModerateUserFormMutation>(
            RequestModerateUserFormDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "RequestModerateUserForm",
        "mutation"
      )
    },
    ApproveModerateUserForm(
      variables: ApproveModerateUserFormMutationVariables,
      requestHeaders?: Dom.RequestInit["headers"]
    ): Promise<{
      data: ApproveModerateUserFormMutation
      extensions?: any
      headers: Dom.Headers
      status: number
    }> {
      return withWrapper(
        (wrappedRequestHeaders) =>
          client.rawRequest<ApproveModerateUserFormMutation>(
            ApproveModerateUserFormDocumentString,
            variables,
            { ...requestHeaders, ...wrappedRequestHeaders }
          ),
        "ApproveModerateUserForm",
        "mutation"
      )
    },
  }
}
export type Sdk = ReturnType<typeof getSdk>
