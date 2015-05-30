package cli
// put bundle under test by inputting yaml from string and checking
// that output Dockerfile is expected
//
// would be nice to run the container and curl it, but that is also
// checked manually
//
// we manually check that output Dockerfile can be built
//
// so we're just checking that *strings* match, but not that the
// functionality is as expected
//
// we probably shouldn't test the wrappers anyway, since they should be
// tested in a separate repo (so we can have language agnosticism)
